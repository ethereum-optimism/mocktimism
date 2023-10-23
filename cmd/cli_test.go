package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
)

func TestCliConfigCommand(t *testing.T) {
	// Create a temp file to act as the config
	tmpfile, err := os.CreateTemp("", "test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
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

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
	expectedConfig := config.Config{
		Profiles: map[string]config.Profile{
			"default": {
				State:  "/path/to/state",
				Silent: false,
				Chains: []config.Chain{
					{
						ID:                 "mainnet",
						BaseChainID:        "mainnet",
						ForkChainID:        1,
						ForkURL:            "https://mainnet.alchemy.infura.io",
						BlockBaseFeePerGas: 420,
						ChainID:            10,
						GasLimit:           420,
						Accounts:           10,
						Balance:            1000,
						StepsTracing:       true,
						AllowOrigin:        "*",
						Port:               8545,
						Host:               "127.0.0.1",
						BlockTime:          12,
						PruneHistory:       false,
					},
					{
						ID:                 "optimism",
						BaseChainID:        "mainnet",
						ForkChainID:        10,
						ForkURL:            "https://op.alchemy.infura.io",
						BlockBaseFeePerGas: 420,
						ChainID:            10,
						GasLimit:           420,
						Accounts:           10,
						Balance:            1000,
						StepsTracing:       true,
						AllowOrigin:        "*",
						Port:               8546,
						Host:               "127.0.0.1",
						BlockTime:          2,
						PruneHistory:       false,
					},
				},
			},
		},
	}
	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	app := newCli("testCommit", "testDate")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = app.Run([]string{"appName", "config", "--config", tmpfile.Name()})
	require.NoError(t, err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	expectedBytes, _ := toml.Marshal(expectedConfig)
	require.Equal(t, string(expectedBytes), string(out))
}

func TestCliConfigCommandJson(t *testing.T) {
	// Create a temp file to act as the config
	tmpfile, err := os.CreateTemp("", "test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
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

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
	expectedConfig := config.Config{
		Profiles: map[string]config.Profile{
			"default": {
				State:  "/path/to/state",
				Silent: false,
				Chains: []config.Chain{
					{
						ID:                 "mainnet",
						BaseChainID:        "mainnet",
						ForkChainID:        1,
						ForkURL:            "https://mainnet.alchemy.infura.io",
						BlockBaseFeePerGas: 420,
						ChainID:            10,
						GasLimit:           420,
						Accounts:           10,
						Balance:            1000,
						StepsTracing:       true,
						AllowOrigin:        "*",
						Port:               8545,
						Host:               "127.0.0.1",
						BlockTime:          12,
						PruneHistory:       false,
					},
					{
						ID:                 "optimism",
						BaseChainID:        "mainnet",
						ForkChainID:        10,
						ForkURL:            "https://op.alchemy.infura.io",
						BlockBaseFeePerGas: 420,
						ChainID:            10,
						GasLimit:           420,
						Accounts:           10,
						Balance:            1000,
						StepsTracing:       true,
						AllowOrigin:        "*",
						Port:               8546,
						Host:               "127.0.0.1",
						BlockTime:          2,
						PruneHistory:       false,
					},
				},
			},
		},
	}
	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	app := newCli("testCommit", "testDate")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = app.Run([]string{"appName", "config", "--config", tmpfile.Name(), "--json"})
	require.NoError(t, err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	expectedBytes, _ := json.MarshalIndent(expectedConfig, "", "\t")
	require.Equal(t, string(expectedBytes), string(out))
}
