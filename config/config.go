package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

type Config struct {
	Profiles map[string]Profile `toml:"profile"`
}

type Profile struct {
	State  string  `toml:"state"`
	Silent bool    `toml:"silent"`
	Chains []Chain `toml:"chains"`
}

type Chain struct {
	ID                 string `toml:"id"`
	BaseChainID        string `toml:"base_chain_id"`
	ForkChainID        int64  `toml:"fork_chain_id"`
	ForkURL            string `toml:"fork_url"`
	BlockBaseFeePerGas int64  `toml:"block_base_fee_per_gas"`
	ChainID            int64  `toml:"chain_id"`
	GasLimit           int64  `toml:"gas_limit"`
	Accounts           int    `toml:"accounts"`
	Balance            int    `toml:"balance"`
	StepsTracing       bool   `toml:"steps-tracing"`
	AllowOrigin        string `toml:"allow-origin"`
	Port               int    `toml:"port"`
	Host               string `toml:"host"`
	BlockTime          int    `toml:"block_time"`
	PruneHistory       bool   `toml:"prune_history"`
}

func LoadNewConfig(log log.Logger, path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	data = []byte(os.ExpandEnv(string(data)))
	log.Debug("parsed new config file", "data", string(data))

	var md toml.MetaData
	md, err = toml.Decode(string(data), &cfg)
	if err != nil {
		log.Error("failed to decode new config file", "err", err)
		return cfg, err
	}

	if len(md.Undecoded()) > 0 {
		log.Error("unknown fields in new config file", "fields", md.Undecoded())
		return cfg, err
	}

	log.Debug("loaded new configuration", "config", cfg)

	return cfg, nil
	// TODO add more validation checks https://github.com/ethereum-optimism/mocktimism/issues/2
}
