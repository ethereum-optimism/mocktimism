package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	// The mocktimism name of the chain.
	Name string `toml:"name"`
	// Set the chain id
	ChainID uint `toml:"chain_id"`
	// The base chain ID when the chain is a rollup
	// If set to 0 or the current chain's ID, the chain is considered an L1 chain
	BaseChainID uint `toml:"base_chain_id"`
	// Specify chain id to skip fetching it from remote endpoint. This enables offline-start mode.
	// You still must pass both `--fork-url` and `--fork-block-number`, and already have your required state cached on disk, anything missing locally would be fetched
	// from the remote.
	ForkChainID uint `toml:"fork_chain_id"`
	// Fetch state over a remote endpoint instead of starting from an empty state.
	// If you want to fetch state from a specific block number, add a block number like `http://localhost:8545@1400000` or use the `fork-block-number` option.
	ForkURL string `toml:"fork_url"`
	// The base fee in a block
	BlockBaseFeePerGas uint `toml:"block_base_fee_per_gas"`
	// Set the gas limit
	GasLimit uint `toml:"gas_limit"`
	// The number of accounts to pre-fund
	Accounts uint `toml:"accounts"`
	// The initial balance of each account
	Balance uint `toml:"balance"`
	// Enable steps tracing used for debug calls returning geth-style traces
	StepsTracing bool `toml:"steps-tracing"`
	//  Set the CORS allow_origin
	AllowOrigin string `toml:"allow-origin"`
	// The port the server will listen on
	Port uint `toml:"port"`
	// The host the server will listen on
	Host string `toml:"host"`
	// Block time in seconds for interval mining.
	BlockTime uint `toml:"block_time"`
	//  Don't keep full chain history. If a number argument is specified, at most this number of states is kept in memory.
	PruneHistory uint `toml:"prune_history"`
}

var DefaultProfile = Profile{
	State:  "",
	Silent: false,
	Chains: []Chain{
		{
			Name:               "L1",
			BaseChainID:        900,
			ForkChainID:        0,
			ForkURL:            "",
			BlockBaseFeePerGas: 1000000000,
			ChainID:            900,
			GasLimit:           30_000_000,
			Accounts:           10,
			Balance:            1000000000000000000,
			StepsTracing:       false,
			AllowOrigin:        "*",
			Port:               8545,
			Host:               "127.0.0.1",
			BlockTime:          0,
			PruneHistory:       0,
		},
		{
			Name:               "L2",
			BaseChainID:        900,
			ForkChainID:        0,
			ForkURL:            "",
			BlockBaseFeePerGas: 1000000000,
			ChainID:            901,
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

func validateChains(chains []Chain) ([]Chain, []error) {
	var errs []error
	chainIDs := make(map[uint]bool)
	forkURLs := make(map[string]bool)
	ports := make(map[uint]bool)

	for i, chain := range chains {
		// Validate uniqueness of ChainID and ForkChainID
		if chainIDs[chain.ChainID] || chainIDs[chain.ForkChainID] {
			errs = append(errs, fmt.Errorf("duplicate ChainID or ForkChainID detected for chain: %s", chain.Name))
		}
		if ports[chain.Port] {
			errs = append(errs, fmt.Errorf("duplicate port detected for chain: %s", chain.Name))
		}

		// Validate BaseChainID
		if chain.BaseChainID != 0 && chain.BaseChainID != chain.ChainID {
			l1Exists := false
			for _, c := range chains {
				if c.ChainID == chain.BaseChainID || c.ForkChainID == chain.BaseChainID {
					l1Exists = true
					break
				}
			}
			if !l1Exists {
				errs = append(errs, fmt.Errorf("no matching L1 BaseChainID found for L2 chain: %s", chain.Name))
			}
		}

		// Validate ForkURL conditions.
		if chain.ChainID != 0 && chain.ForkURL != "" && chain.ChainID != chain.ForkChainID {
			errs = append(errs, fmt.Errorf("cannot set both ChainID and ForkURL for chain: %s. Did you mean to set ForkChainID?", chain.Name))
		}
		if chain.ForkChainID != 0 && chain.ForkURL == "" {
			errs = append(errs, fmt.Errorf("ForkURL must be set if ForkChainID is provided for chain: %s", chain.Name))
		}
		if chain.ForkURL != "" && forkURLs[chain.ForkURL] {
			errs = append(errs, fmt.Errorf("duplicate ForkURL detected: %s", chain.ForkURL))
		}
		forkURLs[chain.ForkURL] = true

		// Defaults
		if chain.Host == "" {
			chain.Host = "127.0.0.1"
		}
		if chain.ChainID == 0 && chain.ForkURL == "" {
			chain.ChainID = findAvailableChainID(chainIDs, 900)
		}
		if chain.Port == 0 {
			chain.Port = findAvailablePort(ports, 8545)
		}

		if chain.Name == "" {
			chain.Name = fmt.Sprintf("%d", chain.ChainID)
		}

		chains[i] = chain
		if chain.ChainID != 0 {
			chainIDs[chain.ChainID] = true
		}
		if chain.ForkChainID != 0 {
			chainIDs[chain.ForkChainID] = true
		}
		ports[chain.Port] = true
	}

	return chains, errs
}

func findAvailableChainID(chainIDs map[uint]bool, startID uint) uint {
	for {
		if !chainIDs[startID] {
			return startID
		}
		startID++
	}
}
func findAvailablePort(ports map[uint]bool, startPort uint) uint {
	for {
		if !ports[startPort] {
			return startPort
		}
		startPort++
	}
}

func validateProfile(profile Profile, path string) (Profile, []error) {
	if profile.State == "" {
		profile.State = DefaultProfile.State
	}
	if !profile.Silent {
		profile.Silent = DefaultProfile.Silent
	}
	if len(profile.Chains) == 0 {
		profile.Chains = DefaultProfile.Chains
	}

	// Keep state "" if it is not set
	// otherwise normalize to an absolute path
	if profile.State != "" {
		// normalize state to absolute path
		baseDir := filepath.Dir(path)
		profile.State = filepath.Clean(filepath.Join(baseDir, profile.State))
	}

	validatedChains, errs := validateChains(profile.Chains)

	profile.Chains = validatedChains

	return profile, errs
}

func LoadNewConfig(log log.Logger, path string) (Config, []error) {
	errs := []error{}
	if path == "" {
		return Config{
			Profiles: map[string]Profile{
				"default": DefaultProfile,
			},
		}, errs
	}
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, append(errs, err)
	}

	data = []byte(os.ExpandEnv(string(data)))
	log.Debug("parsed new config file", "data", string(data))

	var md toml.MetaData
	md, err = toml.Decode(string(data), &cfg)
	if err != nil {
		log.Error("failed to decode new config file", "err", err)
		return cfg, append(errs, err)
	}

	if len(md.Undecoded()) > 0 {
		log.Error("unknown fields in new config file", "fields", md.Undecoded())
		errs = append(errs, fmt.Errorf("unknown fields in new config file: %v", md.Undecoded()))
	}

	log.Debug("loaded new configuration", "config", cfg)

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{
			"default": DefaultProfile,
		}
	}

	if len(cfg.Profiles) == 0 {
		errs = append(errs, fmt.Errorf("no profiles found in config file"))
	}

	for profileName, profile := range cfg.Profiles {
		profileWithDefaults, profileErrs := validateProfile(profile, path)
		errs = append(errs, profileErrs...)
		cfg.Profiles[profileName] = profileWithDefaults
	}

	return cfg, errs
}
