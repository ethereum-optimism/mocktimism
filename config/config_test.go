package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "/path/to/state"
silent = false

# l1 chain
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1

# Fork options
fork_chain_id = 1
fork_url = "https://mainnet.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
gas_limit = 420

# EVM options
accounts = 10
balance = 1000
steps-tracing = true

# Server options
allow-origin = "*"
port = 8545
host = "127.0.0.1"
block_time = 12
prune_history = 0

# l2 chain
[[profile.default.chains]]
name = "optimism"
base_chain_id = 1

# Fork options
fork_chain_id = 10
fork_url = "https://op.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
gas_limit = 420

# EVM options
accounts = 10
balance = 1000
steps-tracing = true

# Server options
allow-origin = "*"
port = 8546
host = "127.0.0.1"
block_time = 2
prune_history = 0
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	err = tmpfile.Close()
	require.NoError(t, err)

	logger := testlog.Logger(t, log.LvlInfo)

	conf, err := LoadNewConfig(logger, tmpfile.Name())

	require.NoError(t, err)

	for profileName, config := range conf.Profiles {
		require.Equal(t, "default", profileName)
		// Global
		baseDir := filepath.Dir(tmpfile.Name())
		expectedState := filepath.Join(baseDir, "/path/to/state")
		require.Equal(t, expectedState, config.State)
		require.Equal(t, false, config.Silent)
		// First chain (L1 mainnet)
		require.Len(t, config.Chains, 2) // Ensure we have 2 chain configurations

		chain1 := config.Chains[0]
		require.Equal(t, "mainnet", chain1.Name)
		require.Equal(t, uint(1), chain1.BaseChainID)
		// Fork options for the first chain
		require.Equal(t, uint(1), chain1.ForkChainID)
		require.Equal(t, "https://mainnet.alchemy.infura.io", chain1.ForkURL)
		require.Equal(t, uint(420), chain1.BlockBaseFeePerGas)
		// Chain options for the first chain
		require.Equal(t, uint(0), chain1.ChainID)
		require.Equal(t, uint(420), chain1.GasLimit)
		// EVM options for the first chain
		require.Equal(t, uint(10), chain1.Accounts)
		require.Equal(t, uint(1000), chain1.Balance)
		require.True(t, chain1.StepsTracing)
		// Server options for the first chain
		require.Equal(t, "*", chain1.AllowOrigin)
		require.Equal(t, uint(8545), chain1.Port)
		require.Equal(t, "127.0.0.1", chain1.Host)
		require.Equal(t, uint(12), chain1.BlockTime)
		require.Equal(t, uint(0), chain1.PruneHistory)

		// Second chain (L2 optimism)
		chain2 := config.Chains[1]
		require.Equal(t, "optimism", chain2.Name)
		require.Equal(t, uint(1), chain2.BaseChainID)
		// Fork options for the second chain
		require.Equal(t, uint(10), chain2.ForkChainID)
		require.Equal(t, "https://op.alchemy.infura.io", chain2.ForkURL)
		require.Equal(t, uint(420), chain2.BlockBaseFeePerGas)
		// Chain options for the second chain
		require.Equal(t, uint(0), chain2.ChainID)
		require.Equal(t, uint(420), chain2.GasLimit)
		// EVM options for the second chain
		require.Equal(t, uint(10), chain2.Accounts)
		require.Equal(t, uint(1000), chain2.Balance)
		require.True(t, chain2.StepsTracing)
		// Server options for the second chain
		require.Equal(t, "*", chain2.AllowOrigin)
		require.Equal(t, uint(8546), chain2.Port)
		require.Equal(t, "127.0.0.1", chain2.Host)
		require.Equal(t, uint(2), chain2.BlockTime)
		require.Equal(t, uint(0), chain2.PruneHistory)
	}
}

func TestLoadConfigDefaultsNoToml(t *testing.T) {
	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	conf, err := LoadNewConfig(logger, "")
	require.NoError(t, err)

	// Now, let's verify that the defaults are set
	for profileName, config := range conf.Profiles {
		require.Equal(t, "default", profileName)
		require.Equal(t, config, DefaultProfile)
	}
}

func TestLoadConfigDefaultsWithEmptyToml(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := ``

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	conf, err := LoadNewConfig(logger, tmpfile.Name())

	require.NoError(t, err)

	// Now, let's verify that the defaults are set
	for profileName, config := range conf.Profiles {
		require.Equal(t, "default", profileName)

		require.Equal(t, DefaultProfile, config)
	}
}

func TestValidatesUniqChainIds(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
chain_id = 1
[[profile.default.chains]]
name = "optimism"
base_chain_id = 1
chain_id = 10
[[profile.default.chains]]
name = "optimism"
base_chain_id = 1
chain_id = 10
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	_, err = LoadNewConfig(logger, tmpfile.Name())
	require.Error(t, err, "duplicate ChainID or ForkChainID detected for chain: %s")
}

func TestValidatesL1Exists(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
chain_id = 1
[[profile.default.chains]]
name = "optimism"
# this should be 1 not 2
base_chain_id = 2
chain_id = 10
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	_, err = LoadNewConfig(logger, tmpfile.Name())
	require.Error(t, err, "no matching L1 BaseChainID found for L2 chain:")
}

func TestValidatesChainIdAndForkUrl(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
chain_id = 1
[[profile.default.chains]]
name = "optimism"
# this should be 1 not 2
base_chain_id = 2
chain_id = 10
fork_url = "https://op.alchemy.infura.io"
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	_, err = LoadNewConfig(logger, tmpfile.Name())
	require.Error(t, err)
}

func TestDefaultsPortsAndHost(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
chain_id = 1
[[profile.default.chains]]
name = "optimism"
base_chain_id = 1
chain_id = 10
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	cfg, err := LoadNewConfig(logger, tmpfile.Name())
	require.NoError(t, err)

	for _, profile := range cfg.Profiles {
		for i, chain := range profile.Chains {
			expectedPort := DefaultProfile.Chains[0].Port + uint(i)
			require.Equal(t, DefaultProfile.Chains[0].Host, chain.Host)
			require.Equal(t, expectedPort, chain.Port)
		}
	}
}

func TestForkBlockNumber(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
fork_chain_id = 1
fork_block_number = 1234
fork_url = "https://mainnet.alchemy.infura.io"
[[profile.default.chains]]
name = "optimism"
# this should be 1 not 2
base_chain_id = 1
fork_chain_id = 10
fork_url = "https://op.alchemy.infura.io"
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	cfg, errs := LoadNewConfig(logger, tmpfile.Name())
	err = errors.Join(errs...)
	require.NoError(t, err)
	for _, profile := range cfg.Profiles {
		for _, chain := range profile.Chains {
			if chain.Name == "mainnet" {
				require.Equal(t, uint(1234), chain.ForkBlockNumber)
			} else if chain.Name == "optimism" {
				require.Equal(t, uint(0), chain.ForkBlockNumber)
			} else {
				t.Errorf("unexpected chain name: %s", chain.Name)
			}
		}
	}
}

func TestForkBlockNumberWithNoForkUrlError(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
fork_chain_id = 1
fork_block_number = 1234
[[profile.default.chains]]
name = "optimism"
# this should be 1 not 2
base_chain_id = 2
chain_id = 10
fork_url = "https://op.alchemy.infura.io"
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	_, errs := LoadNewConfig(logger, tmpfile.Name())
	err = errors.Join(errs...)
	require.Error(t, err, "ForkBlockNumber is set but no ForkURL is not provided for chain: optimism")
}

func TestForkBlockNumberOnL2Error(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "default_test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
[profile.default]
state = "path/to/state"
[[profile.default.chains]]
name = "mainnet"
base_chain_id = 1
fork_chain_id = 1
fork_block_number = 1234
[[profile.default.chains]]
name = "optimism"
# this should be 1 not 2
base_chain_id = 2
chain_id = 10
fork_url = "https://op.alchemy.infura.io"
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	// Load the configuration
	logger := testlog.Logger(t, log.LvlInfo)
	_, errs := LoadNewConfig(logger, tmpfile.Name())
	err = errors.Join(errs...)
	require.Error(t, err, "ForkBlockNumber cannot be set for L2 network: optimism. Try setting fork-block-number on the L1 network instead")
}
