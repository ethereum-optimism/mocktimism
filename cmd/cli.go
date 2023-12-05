package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/orchestrator"
	"github.com/ethereum-optimism/mocktimism/services/anvil"
	"github.com/ethereum-optimism/mocktimism/services/geth"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

func newCli(GitCommit string, GitDate string) *cli.App {
	configFlags := []cli.Flag{
		ConfigFlag,
		JsonFlag,
	}
	configFlags = append(configFlags, oplog.CLIFlags("MOCKTIMISM")...)
	return &cli.App{
		Version:              params.VersionWithCommit(GitCommit, GitDate),
		Description:          "A cli wrapper around anvil for spinning up devnets",
		EnableBashCompletion: true,
		Action: func(ctx *cli.Context) error {
			o := orchestrator.NewOrchestrator("mocktimism")
			// TODO extract the code to run every process into a reusable function https://github.com/ethereum-optimism/mocktimism/issues/73
			var wg sync.WaitGroup
			errCh := make(chan error, 5)

			log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
			oplog.SetGlobalLogHandler(log.GetHandler())
			cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
			if err != nil {
				log.Error("failed to load config", "err", err)
				return err
			}

			processCtx, processCancel := context.WithCancel(ctx.Context)

			runService := func(start func(ctx context.Context) error) {
				wg.Add(1)
				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Error("Mocktimism had an unexpected fatal error", "err", err)
							debug.PrintStack()
							errCh <- fmt.Errorf("panic: %v", err)
						}

						processCancel()
						wg.Done()
					}()

					errCh <- start(processCtx)
				}()
			}

			for _, profile := range cfg.Profiles {
				for _, chain := range profile.Chains {
					// Start anvil
					anvil, err := anvil.NewAnvilService(chain.Name, log, chain)
					o.Register(anvil)
					if err != nil {
						log.Error("failed to create anvil service", "err", err)
						return err
					}
					runService(anvil.Start)
				}
			}
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:        "config",
				Flags:       configFlags,
				Description: "Display the current mocktimism config",
				Action:      actionConfig,
			},
			{
				Name:        "",
				Flags:       configFlags,
				Description: "Starts the anvil services",
				Action: func(ctx *cli.Context) error {
					// TODO extract the code to run every process into a reusable function https://github.com/ethereum-optimism/mocktimism/issues/73
					var wg sync.WaitGroup
					errCh := make(chan error, 5)

					log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
					oplog.SetGlobalLogHandler(log.GetHandler())
					cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
					if err != nil {
						log.Error("failed to load config", "err", err)
						return err
					}

					processCtx, processCancel := context.WithCancel(ctx.Context)

					runService := func(start func(ctx context.Context) error) {
						wg.Add(1)
						go func() {
							defer func() {
								if err := recover(); err != nil {
									log.Error("Mocktimism had an unexpected fatal error", "err", err)
									debug.PrintStack()
									errCh <- fmt.Errorf("panic: %v", err)
								}

								processCancel()
								wg.Done()
							}()

							errCh <- start(processCtx)
						}()
					}

					for _, profile := range cfg.Profiles {
						for _, chain := range profile.Chains {
							anvil, err := anvil.NewAnvilService(chain.Name, log, chain)
							if err != nil {
								log.Error("failed to create anvil service", "err", err)
								return err
							}
							runService(anvil.Start)
						}
					}
					return nil
				},
			},
			{
				Name:        "geth-l1",
				Flags:       configFlags,
				Description: "Starts the geth l1 services",
				Action: func(ctx *cli.Context) error {
					log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
					cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
					gc := geth.GethConfig{
						OpGeth:    false,
						HTTPPort:  int(cfg.Profiles["default"].Chains[0].Port),
						Verbosity: 3,
						DataDir:   "/tmp/geth",
					}
					service, err := geth.NewGeth("l1", log, gc, nil)
					if err != nil {
						log.Error("failed to create geth service", "err", err)
						return err
					}
					var wg sync.WaitGroup
					errCh := make(chan error, 5)
					processCtx, processCancel := context.WithCancel(ctx.Context)
					runService := func(start func(ctx context.Context) error) {
						wg.Add(1)
						log.Info("starting new service...")
						go func() {
							log.Info("yes starting new service...")
							defer func() {
								log.Info("Ending")
								if err := recover(); err != nil {
									log.Error("Mocktimism had an unexpected fatal error", "err", err)
									debug.PrintStack()
									errCh <- fmt.Errorf("panic: %v", err)
								}

								processCancel()
								wg.Done()
							}()

							log.Info("Starting service")
							errCh <- start(processCtx)
						}()
					}

					if true {
						service.Start(processCtx)
						ticker := time.NewTicker(2 * time.Second)
						defer ticker.Stop()
						for {
							select {
							case <-ticker.C:
								healthy, err := service.HealthCheck()
								if err != nil || !healthy {
									log.Error("Health check failed:", err)
									processCancel()
									return nil
								}
								log.Info("health check passed")
							case <-ctx.Done():
								return nil
							}
						}
					} else {
						runService(service.Start)
						log.Info("Ending service")
						return nil
					}
				},
			},
			{
				Name:        "anvil",
				Flags:       configFlags,
				Description: "Starts the anvil services",
				Action: func(ctx *cli.Context) error {
					// TODO extract the code to run every process into a reusable function https://github.com/ethereum-optimism/mocktimism/issues/73
					var wg sync.WaitGroup
					errCh := make(chan error, 5)
					processCtx, processCancel := context.WithCancel(ctx.Context)

					log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
					oplog.SetGlobalLogHandler(log.GetHandler())
					cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
					if err != nil {
						log.Error("failed to load config", "err", err)
						return err
					}

					log.Info("Starting services...")
					runService := func(start func(ctx context.Context) error) {
						wg.Add(1)
						log.Info("starting new service...")
						go func() {
							log.Info("yes starting new service...")
							defer func() {
								log.Info("Ending")
								if err := recover(); err != nil {
									log.Error("Mocktimism had an unexpected fatal error", "err", err)
									debug.PrintStack()
									errCh <- fmt.Errorf("panic: %v", err)
								}

								processCancel()
								wg.Done()
							}()

							log.Info("Starting service")
							errCh <- start(processCtx)
						}()
					}

					for _, profile := range cfg.Profiles {
						// TODO we only want to use default or the specified profile
						for i, chain := range profile.Chains {
							anvil, err := anvil.NewAnvilService(chain.Name, log, chain)
							if err != nil {
								log.Error("failed to create anvil service", "err", err)
								return err
							}
							log.Info("Starting chain", "chain", chain.Name, "index", i)
							runService(anvil.Start)
						}
					}
					wg.Wait()
					return nil
				},
			},
		},
	}
}
