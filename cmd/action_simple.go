// This is just a simple devnet that doesn't read from config at all
// this will be deleted eventually it's only here for development purposes
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/ethereum-optimism/mocktimism/config"

	"github.com/ethereum-optimism/mocktimism/prysm"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core"
	"github.com/pelletier/go-toml"
	"github.com/urfave/cli/v2"
)

//go:embed static/geth/genesis.json
var genesisFile []byte

//go:embed static/keystore/UTC--2022-08-19T17-38-31.257380510Z--123463a4b065722e99115d6c222f267d9cabb524
var keystoreFile []byte

func ActionSimple(ctx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "mocktimism")
	oplog.SetGlobalLogHandler(log.GetHandler())

	cfg, err := config.LoadNewConfig(log, ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("failed to load config", "errors", err)
		return err
	}

	var genesis core.Genesis
	if err := json.Unmarshal(genesisFile, &genesis); err != nil {
		log.Error("failed to unmarshal genesis file", "error", err)
		return err
	}

	ks := keystore.NewKeyStore("/tmp/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	key, err := ks.Import(keystoreFile, "", "") // You might need to replace the second and third arguments with the actual passphrase
	if err != nil {
		log.Error("failed to import keystore file", "error", err)
		return err
	}

	err = prysmmanager.StartBeaconChain(genesis)
	if err != nil {
		log.Error("failed to start beacon chain", "error", err)
		return err
	}

	return nil
}
