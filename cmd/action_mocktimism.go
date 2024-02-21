package main

import (
	"context"

	"github.com/ethereum-optimism/mocktimism"
	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/urfave/cli/v2"
)

func actionMocktimism(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	logger := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	oplog.SetGlobalLogHandler(logger.Handler())
	logger.Debug("running mocktimism...")

	cfg, err := config.LoadNewConfig(logger, ctx.String(ConfigFlag.Name))
	if err != nil {
		logger.Error("failed to load config", "errors", err)
		return nil, err
	}
	return mocktimism.NewMocktimism(ctx.Context, logger, shutdown, &cfg)
}
