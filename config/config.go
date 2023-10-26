package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/log"
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
	// The mocktimism name of the chain
	ID string `toml:"id"`
	// The base chain ID when the chain is a rollup
	// If set to 0 or the current chain's ID, the chain is considered an L1 chain
	BaseChainID int64 `toml:"base_chain_id"`
	// Specify chain id to skip fetching it from remote endpoint. This enables offline-start mode.
	// You still must pass both `--fork-url` and `--fork-block-number`, and already have your required state cached on disk, anything missing locally would be fetched
	// from the remote.
	ForkChainID int64 `toml:"fork_chain_id"`
	// Fetch state over a remote endpoint instead of starting from an empty state.
	// If you want to fetch state from a specific block number, add a block number like `http://localhost:8545@1400000` or use the `fork-block-number` option.
	ForkURL string `toml:"fork_url"`
	// The base fee in a block
	BlockBaseFeePerGas int64 `toml:"block_base_fee_per_gas"`
	// Set the chain id
	ChainID int64 `toml:"chain_id"`
	// Set the gas limit
	GasLimit int64 `toml:"gas_limit"`
	// The number of accounts to pre-fund
	Accounts int `toml:"accounts"`
	// The initial balance of each account
	Balance int `toml:"balance"`
	// Enable steps tracing used for debug calls returning geth-style traces
	StepsTracing bool `toml:"steps-tracing"`
	//  Set the CORS allow_origin
	AllowOrigin string `toml:"allow-origin"`
	// The port the server will listen on
	Port int `toml:"port"`
	// The host the server will listen on
	Host string `toml:"host"`
	// Block time in seconds for interval mining.
	BlockTime int `toml:"block_time"`
	//  Don't keep full chain history. If a number argument is specified, at most this number of states is kept in memory.
	PruneHistory int `toml:"prune_history"`
}

var DefaultProfile = Profile{
	State:  "",
	Silent: true,
	Chains: []Chain{
		{
			ID:                 "L1",
			BaseChainID:        900,
			ForkChainID:        900,
			ForkURL:            "",
			BlockBaseFeePerGas: 1000000000,
			ChainID:            900,
			GasLimit:           30_000_000,
			Accounts:           10,
			Balance:            1000000000000000000,
			StepsTracing:       false,
			AllowOrigin:        "*",
			Port:               8545,
			Host:               "localhost",
			BlockTime:          0,
			PruneHistory:       0,
		},
		{
			ID:                 "L2",
			BaseChainID:        901,
			ForkChainID:        901,
			ForkURL:            "",
			BlockBaseFeePerGas: 1000000000,
			ChainID:            900,
			GasLimit:           30_000_000,
			Accounts:           10,
			Balance:            1000000000000000000,
			StepsTracing:       false,
			AllowOrigin:        "*",
			Port:               9545,
			Host:               "localhost",
			BlockTime:          0,
			PruneHistory:       0,
		},
	},
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
