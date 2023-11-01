package anvil

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/stretchr/testify/require"
)

func TestGetL2BlockNumberAtSpecificL1Block(t *testing.T) {
	logger := log.New("module", "test")
	cfg := config.Chain{
		Host: "127.0.0.1",
		Port: 8545,
	}

	// Initialize the AnvilService
	service, err := NewAnvilService("TestService", logger, cfg)
	require.NoError(t, err)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(900))
	require.NoError(t, err)
	auth.GasLimit = 8000000
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{from: {Balance: big.NewInt(params.Ether)}}, 50_000_000)
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	defer backend.Close()

	l2OutputOracleAddress, _, _, err := bindings.DeployL2OutputOracle(
		opts,
		backend,
		big.NewInt(10),
		big.NewInt(2),
		big.NewInt(100),
	)
	require.NoError(t, err)

	backend.Commit()

	l1BlockNumber := big.NewInt(1) // As an example
	_, err = service.L2BlockNumber(backend, l2OutputOracleAddress, l1BlockNumber)
	require.NoError(t, err)

	// TODO add proposals

	// expectedL2BlockNumber := big.NewInt(42)
	// require.Equal(t, expectedL2BlockNumber, l2BlockNumber)
}
