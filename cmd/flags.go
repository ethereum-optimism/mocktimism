package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func findConfigFile() string {
	currentDir, _ := os.Getwd()
	return findConfigRecursively(currentDir, "mocktimism.toml")
}

func findConfigRecursively(dir, fileName string) string {
	configPath := filepath.Join(dir, fileName)
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		// We've reached the root directory, and the file is not found.
		return ""
	}

	return findConfigRecursively(parentDir, fileName)
}

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Value:   findConfigFile(),
		Aliases: []string{"c"},
		Usage:   "path to config file",
		EnvVars: []string{"MOCKTIMISM_CONFIG"},
	}
	JsonFlag = &cli.BoolFlag{
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "print config in JSON form",
		EnvVars: []string{"MOCKTIMISM_CONFIG_JSON"},
	}
)
