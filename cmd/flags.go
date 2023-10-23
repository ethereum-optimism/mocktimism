package main

import (
	"github.com/urfave/cli/v2"
)

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Value:   "./mocktimism.toml",
		Aliases: []string{"c"},
		Usage:   "path to config file",
		EnvVars: []string{"MOCKTIMISM_CONFIG"},
	}
	JsonFlag = &cli.BoolFlag{
		Name: 	 "json",
		Aliases: []string{"j"},
		Usage:	 "print config in JSON form",
		EnvVars: []string{"MOCKTIMISM_CONFIG_JSON"},
	}
)
