package main

import (
	"context"

	"github.com/ethereum-optimism/mocktimism"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

func runMocktimism(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	logger := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	// TODO why is this not working with logger.Handler() like in https://github.com/ethereum-optimism/optimism/blob/develop/indexer/cmd/indexer/cli.go#L48
	// For some reason Handler does not exist
	// oplog.SetGlobalLogHandler(logger.Handler())
	logger.Debug("running mocktimism...")

	return mocktimism.NewMocktimism(ctx.Context, logger, shutdown)
}
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
		Action:               cliapp.LifecycleCmd(runMocktimism),
		Commands: []*cli.Command{
			{
				Name:        "config",
				Flags:       configFlags,
				Description: "Display the current mocktimism config",
				Action:      actionConfig,
			},
		},
	}
}
