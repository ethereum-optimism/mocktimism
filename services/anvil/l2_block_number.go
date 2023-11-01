package anvil

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// TODO connect this to superchain repo https://github.com/ethereum-optimism/mocktimism/issues/79
func (a *AnvilService) L2BlockNumber(l1Client bind.ContractBackend, l2OutputOracleAddr common.Address, l1BlockNumber *big.Int) (*big.Int, error) {
	l2OutputOracle, err := bindings.NewL2OutputOracle(l2OutputOracleAddr, l1Client)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{
		Context:     context.Background(),
		BlockNumber: l1BlockNumber,
	}

	l2Height, err := l2OutputOracle.LatestBlockNumber(opts)
	if err != nil {
		return nil, err
	}

	return l2Height, nil
}
