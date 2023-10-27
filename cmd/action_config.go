package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum-optimism/mocktimism/config"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/pelletier/go-toml"
	"github.com/urfave/cli/v2"
)

func actionConfig(ctx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	oplog.SetGlobalLogHandler(log.GetHandler())
	cfg, errs := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
	if len(errs) > 0 {
		err := errors.Join(errs...)
		log.Error("failed to load config", "errors", err)
		return err
	}
	if ctx.Bool(JsonFlag.Name) {
		s, _ := json.MarshalIndent(cfg, "", "\t")
		fmt.Print(string(s))
	} else {
		s, _ := toml.Marshal(cfg)
		fmt.Print(string(s))
	}
	return nil
}
