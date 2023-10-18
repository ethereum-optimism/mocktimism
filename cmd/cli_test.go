package main

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestCliConfigCommand(t *testing.T) {
	// Create a temp file to act as the config
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

	app := newCli("testCommit", "testDate")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = app.Run([]string{"appName", "config", "--config", tmpfile.Name()})
	require.NoError(t, err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	expected := `
t=2023-10-18T09:11:41-0700 lvl=eror msg=\"unknown fields in new config file\" role=mocktimism fields=\"[profile.default profile.default.state profile.default.silent profile.default.chains profile.default.chains.id profile.default.chains.base_chain_id profile.default.chains.fork_chain_id profile.default.chains.fork_url profile.default.chains.block_base_fee_per_gas profile.default.chains.chain_id profile.default.chains.gas_limit profile.default.chains.accounts profile.default.chains.balance profile.default.chains.steps-tracing profile.default.chains.allow-origin profile.default.chains.port profile.default.chains.host profile.default.chains.block_time profile.default.chains.prune_history profile.default.chains profile.default.chains.id profile.default.chains.base_chain_id profile.default.chains.fork_chain_id profile.default.chains.fork_url profile.default.chains.block_base_fee_per_gas profile.default.chains.chain_id profile.default.chains.gas_limit profile.default.chains.accounts profile.default.chains.balance profile.default.chains.steps-tracing profile.default.chains.allow-origin profile.default.chains.port profile.default.chains.host profile.default.chains.block_time profile.default.chains.prune_history]\"\n{\n\t\"Profiles\": null\n}
  `
	require.Equal(t, expected, string(out))
}
