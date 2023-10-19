package main

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum-optimism/mocktimism/config"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/urfave/cli/v2"
)

func actionConfig(ctx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	oplog.SetGlobalLogHandler(log.GetHandler())
	cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("failed to load config", "err", err)
		return err
	}
	s, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Print(string(s))
	return nil
}
