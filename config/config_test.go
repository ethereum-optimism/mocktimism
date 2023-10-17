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
		[global]
		state = '/path/to/my/persisted-state'
		silent = false

		[l1]
		id = 'mainnet'

		[[fork]]
		fork_chain_id = 1
		fork_url = "https://mainnet.alchemy.infura"
		block_base_fee_per_gas = 420

		[[environment]]
		block_base_fee_per_gas = 420
		chain_id = 10
		gas_limit = 420

		[[evm]]
		accounts = 10
		balance = 1000
		steps-tracing = true

		[[server]]
		allow-origin = "*"
		port = 8545
		host = "127.0.0.1"
		block_time = 12
		prune_history = false

		[l2]
		id = "optimism"
		l1 = "mainnet"

		[[fork]]
		fork_chain_id = 10
		fork_url = "https://op.alchemy.infura"
		block_base_fee_per_gas = 420

		[[environment]]
		l1_block_base_per_gas = 420
		block_base_fee_per_gas = 420
		chain_id = 10
		gas_limit = 420

		[[evm]]
		accounts = 10
		balance = 1000
		steps-tracing = true

		[[server]]
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

	// Global
	require.Equal(t, "/path/to/my/persisted-state", conf.Global.State)
	require.Equal(t, false, conf.Global.Silent)

	// L1
	require.Equal(t, "mainnet", conf.L1.Id)
	require.Equal(t, 1, conf.L1.Fork.ForkChainId)
	require.Equal(t, "https://mainnet.alchemy.infura", conf.L1.Fork.ForkURL)
	require.Equal(t, 420, conf.L1.Environments[0].BlockBaseFeePerGas)
	require.Equal(t, 10, conf.L1.Evms[0].Accounts)
	require.Equal(t, 1000, conf.L1.Evms[0].Balance)
	require.Equal(t, true, conf.L1.Evms[0].StepsTracing)
	require.Equal(t, "*", conf.L1.Servers[0].AllowOrigin)
	require.Equal(t, 8545, conf.L1.Servers[0].Port)
	require.Equal(t, "127.0.0.1", conf.L1.Servers[0].Host)
	require.Equal(t, 12, conf.L1.Servers[0].BlockTime)
	require.Equal(t, false, conf.L1.Servers[0].PruneHistory)

	// L2
	require.Equal(t, "optimism", conf.L2.Id)
	require.Equal(t, 10, conf.L2.Fork.ForkChainId)
	require.Equal(t, "https://op.alchemy.infura", conf.L2.Fork.ForkURL)
	require.Equal(t, 420, *conf.L2.Environments[0].L1BlockBaseFeePerGas)
	require.Equal(t, 10, conf.L2.Evms[0].Accounts)
	require.Equal(t, 1000, conf.L2.Evms[0].Balance)
	require.Equal(t, true, conf.L2.Evms[0].StepsTracing)
	require.Equal(t, "*", conf.L2.Servers[0].AllowOrigin)
	require.Equal(t, 8546, conf.L2.Servers[0].Port)
	require.Equal(t, "127.0.0.1", conf.L2.Servers[0].Host)
	require.Equal(t, 2, conf.L2.Servers[0].BlockTime)
	require.Equal(t, false, conf.L2.Servers[0].PruneHistory)

	// TODO add more validation checks https://github.com/ethereum-optimism/mocktimism/issues/2
}
