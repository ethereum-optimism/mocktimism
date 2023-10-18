package config

import (
	"os"
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
id = "mainnet"
base_chain_id = "mainnet"

# Fork options
fork_chain_id = 1
fork_url = "https://mainnet.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
chain_id = 10
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
prune_history = false

# l2 chain
[[profile.default.chains]]
id = "optimism"
base_chain_id = "mainnet"

# Fork options
fork_chain_id = 10
fork_url = "https://op.alchemy.infura.io"
block_base_fee_per_gas = 420

# Chain options
chain_id = 10
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
prune_history = false
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

	for profileName, config := range conf.Profile {
		require.Equal(t, "default", profileName)
		// Global
		require.Equal(t, "/path/to/state", config.State)
		require.Equal(t, false, config.Silent)
		// First chain (L1 mainnet)
		require.Len(t, config.Chains, 2) // Ensure we have 2 chain configurations

		chain1 := config.Chains[0]
		require.Equal(t, "mainnet", chain1.ID)
		require.Equal(t, "mainnet", chain1.BaseChainID)
		// Fork options for the first chain
		require.Equal(t, int64(1), chain1.ForkChainID)
		require.Equal(t, "https://mainnet.alchemy.infura.io", chain1.ForkURL)
		require.Equal(t, int64(420), chain1.BlockBaseFeePerGas)
		// Chain options for the first chain
		require.Equal(t, int64(10), chain1.ChainID)
		require.Equal(t, int64(420), chain1.GasLimit)
		// EVM options for the first chain
		require.Equal(t, 10, chain1.Accounts)
		require.Equal(t, 1000, chain1.Balance)
		require.True(t, chain1.StepsTracing)
		// Server options for the first chain
		require.Equal(t, "*", chain1.AllowOrigin)
		require.Equal(t, 8545, chain1.Port)
		require.Equal(t, "127.0.0.1", chain1.Host)
		require.Equal(t, 12, chain1.BlockTime)
		require.False(t, chain1.PruneHistory)

		// Second chain (L2 optimism)
		chain2 := config.Chains[1]
		require.Equal(t, "optimism", chain2.ID)
		require.Equal(t, "mainnet", chain2.BaseChainID)
		// Fork options for the second chain
		require.Equal(t, int64(10), chain2.ForkChainID)
		require.Equal(t, "https://op.alchemy.infura.io", chain2.ForkURL)
		require.Equal(t, int64(420), chain2.BlockBaseFeePerGas)
		// Chain options for the second chain
		require.Equal(t, int64(10), chain2.ChainID)
		require.Equal(t, int64(420), chain2.GasLimit)
		// EVM options for the second chain
		require.Equal(t, 10, chain2.Accounts)
		require.Equal(t, 1000, chain2.Balance)
		require.True(t, chain2.StepsTracing)
		// Server options for the second chain
		require.Equal(t, "*", chain2.AllowOrigin)
		require.Equal(t, 8546, chain2.Port)
		require.Equal(t, "127.0.0.1", chain2.Host)
		require.Equal(t, 2, chain2.BlockTime)
		require.False(t, chain2.PruneHistory)
	}

}
