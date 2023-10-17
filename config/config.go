package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/log"
)

type GlobalConfig struct {
	State  string `toml:"state"`
	Silent bool   `toml:"silent"`
}

type ForkConfig struct {
	ForkChainId        int    `toml:"fork_chain_id"`
	ForkURL            string `toml:"fork_url"`
	BlockBaseFeePerGas int    `toml:"block_base_fee_per_gas"`
}

type EnvironmentConfig struct {
	L1BlockBaseFeePerGas *int `toml:"l1_block_base_fee_per_gas,omitempty"`
	BlockBaseFeePerGas   int  `toml:"block_base_fee_per_gas"`
	ChainId              int  `toml:"chain_id"`
	GasLimit             int  `toml:"gas_limit"`
}

type EvmConfig struct {
	Accounts     int  `toml:"accounts"`
	Balance      int  `toml:"balance"`
	StepsTracing bool `toml:"steps-tracing"`
}

type ServerConfig struct {
	AllowOrigin  string `toml:"allow-origin"`
	Port         int    `toml:"port"`
	Host         string `toml:"host"`
	BlockTime    int    `toml:"block_time"`
	PruneHistory bool   `toml:"prune_history"`
}

type ChainConfig struct {
	Id           string              `toml:"id"`
	Fork         ForkConfig          `toml:"fork"`
	Environments []EnvironmentConfig `toml:"environment"`
	Evms         []EvmConfig         `toml:"evm"`
	Servers      []ServerConfig      `toml:"server"`
}

type Configuration struct {
	Global GlobalConfig `toml:"global"`
	L1     ChainConfig  `toml:"l1"`
	L2     ChainConfig  `toml:"l2"`
}

func LoadNewConfig(log log.Logger, path string) (Configuration, error) {
	log.Debug("loading new config", "path", path)

	var cfg Configuration
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	data = []byte(os.ExpandEnv(string(data)))
	log.Debug("parsed new config file", "data", string(data))

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		log.Error("failed to decode new config file", "err", err)
		return cfg, err
	}

	log.Info("loaded new configuration", "config", cfg)
	return cfg, nil

	// TODO add more validation checks https://github.com/ethereum-optimism/mocktimism/issues/2
}
