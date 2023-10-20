package main

import (
	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/service-discovery"
	"github.com/ethereum-optimism/mocktimism/services/anvil"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/urfave/cli/v2"
)

var (
	devnetName = "mocktimism-devnet"
)

func actionUp(ctx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	oplog.SetGlobalLogHandler(log.GetHandler())
	cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
	log.Debug("Loaded config", "cfg", cfg)

	log.Debug("Starting service discovery registry...")
	servicediscovery.NewServiceDiscovery(devnetName)
	log.Debug("Starting services...")
	log.Debug("Starting l1...")
	anvilService := anvil.NewAnvilService(log.New("role", "l1"))
	if err := anvilService.Start(ctx.Context); err != nil {
		log.Error("failed to start l1", "err", err)
		return err
	}

	if err != nil {
		log.Error("failed to load config", "err", err)
		return err
	}

	return nil
}
