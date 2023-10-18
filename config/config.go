package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

type Config struct {
	Profiles map[string]Profile
}

type Profile struct {
	State  string
	Silent bool
	Chains []Chain
}

type Chain struct {
	ID                 string
	BaseChainID        string
	ForkChainID        int64
	ForkURL            string
	BlockBaseFeePerGas int64
	ChainID            int64
	GasLimit           int64
	Accounts           int
	Balance            int
	StepsTracing       bool
	AllowOrigin        string
	Port               int
	Host               string
	BlockTime          int
	PruneHistory       bool
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

	log.Info("loaded new configuration", "config", cfg)

	return cfg, nil
	// TODO add more validation checks https://github.com/ethereum-optimism/mocktimism/issues/2
}
