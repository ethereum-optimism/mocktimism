package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/services/anvil"
	"github.com/ethereum/go-ethereum/log"
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
name = "mainnet"
base_chain_id = 1

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
prune_history = 0
`

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
	expectedConfig := config.Config{
		Profiles: map[string]config.Profile{
			"default": {
				State:  "/path/to/state",
				Silent: false,
				Chains: []config.Chain{
					{
						Name:               "mainnet",
						BaseChainID:        1,
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
						PruneHistory:       0,
					},
					{
						Name:               "optimism",
						BaseChainID:        1,
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
						PruneHistory:       0,
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
name = "mainnet"
base_chain_id = 1

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
prune_history = 0
`

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
	expectedConfig := config.Config{
		Profiles: map[string]config.Profile{
			"default": {
				State:  "/path/to/state",
				Silent: false,
				Chains: []config.Chain{
					{
						Name:               "mainnet",
						BaseChainID:        1,
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
						PruneHistory:       0,
					},
					{
						Name:               "optimism",
						BaseChainID:        1,
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
						PruneHistory:       0,
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

func TestCliAnvilCommand(t *testing.T) {
	// Create a temp file to act as the config
	tmpfile, err := os.CreateTemp("", "test.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// TODO can move this to shared test fixture folder https://github.com/ethereum-optimism/mocktimism/issues/29
	testData := `
[profile.default]

[[profile.default.chains]]
port = 8545
host = "127.0.0.1"

# l2 chain
[[profile.default.chains]]
port = 9545
host = "127.0.0.1"
`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app := newCli("testCommit", "testDate")

	go func() {
		err := app.RunContext(ctx, []string{"appName", "anvil", "--config", tmpfile.Name()})
		require.NoError(t, err)
	}()

	// Poll for l1 chain healtcheck
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond) // polling every 200ms
	defer ticker.Stop()

	healthy := false
loop:
	for {
		select {
		case <-timeout:
			break loop
		case <-ticker.C:
			// Check l1
			service, err := anvil.NewAnvilService(
				"HealthCheck",
				log.New("module", "test"),
				config.Chain{
					Host: "127.0.0.1",
					Port: 8545,
				})
			if err != nil {
				break loop
			}
			healthy, err = service.HealthCheck()
			if healthy || err != nil {
				break loop
			}
			// Check l2
			service, err = anvil.NewAnvilService(
				"HealthCheck",
				log.New("module", "test"),
				config.Chain{
					Host: "127.0.0.1",
					Port: 9545,
				})
			if err != nil {
				break loop
			}
			healthy, err = service.HealthCheck()
			if healthy || err != nil {
				break loop
			}
		}
	}

	require.NoError(t, err, "Health check failed")
	require.True(t, healthy, "Service is not healthy after waiting for 2 seconds")
}
