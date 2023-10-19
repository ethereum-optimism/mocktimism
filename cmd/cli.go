package main

import (
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

func newCli(GitCommit string, GitDate string) *cli.App {
	flags := []cli.Flag{ConfigFlag}
	flags = append(flags, oplog.CLIFlags("MOCKTIMISM")...)
	return &cli.App{
		Version:              params.VersionWithCommit(GitCommit, GitDate),
		Description:          "A cli wrapper around anvil for spinning up devnets",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "config",
				Flags:       flags,
				Description: "Display the current mocktimism config",
				Action:      actionConfig,
			},
		},
	}
}
