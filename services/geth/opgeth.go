package geth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/op-geth/accounts/abi/bind"
	"github.com/ethereum-optimism/op-geth/common"
	"github.com/ethereum-optimism/op-geth/common/hexutil"
	"github.com/ethereum-optimism/op-geth/core"
	"github.com/ethereum-optimism/op-geth/eth"
	"github.com/ethereum-optimism/op-geth/eth/ethconfig"
	"github.com/ethereum-optimism/op-geth/log"
	"github.com/ethereum-optimism/op-geth/node"
	"github.com/ethereum-optimism/op-geth/p2p"
	"github.com/ethereum-optimism/op-geth/rpc"
)

type GethConfig struct {
	DataDir   string
	Verbosity int
	HTTPPort  int
}

type Geth struct {
	log    log.Logger
	config GethConfig

	node *node.Node
	eth  *eth.Ethereum
}

func NewGeth(logger log.Logger, cfg GethConfig) (*Geth, error) {
	ethCfg := &ethconfig.Config{
		NetworkId: 1,
		Genesis:   core.DefaultGenesisBlock(),
	}

	nodeCfg := &node.Config{
		Name:        "simple-geth",
		HTTPHost:    "0.0.0.0",
		HTTPPort:    cfg.HTTPPort,
		HTTPModules: []string{"web3", "eth", "txpool", "net", "rpc"},
		DataDir:     cfg.DataDir,
		P2P: p2p.Config{
			NoDiscovery: true,
			MaxPeers:    1,
		},
	}

	n, err := node.New(nodeCfg)
	if err != nil {
		return nil, err
	}

	backend, err := eth.New(n, ethCfg)
	if err != nil {
		return nil, err
	}

	return &Geth{
		log:    logger,
		config: cfg,
		node:   n,
		eth:    backend,
	}, nil
}

func (s *Geth) Start(ctx context.Context) error {
	err := s.node.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		s.Close()
	}()

	return nil
}

func (s *Geth) Close() error {
	return s.node.Close()
}

func (s *Geth) HealthCheck() (bool, error) {
	client, err := s.GetClient()
	if client != nil {
		defer client.Close()
	}
	if err != nil {
		return false, err
	}
	_, err = s.BlockNumber(client)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Geth) GetClient() (*rpc.Client, error) {
	httpEndpoint := fmt.Sprintf("http://0.0.0.0:%d", s.config.HTTPPort) // Assuming the node is running locally, hence using 0.0.0.0
	client, err := rpc.Dial(httpEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC: %w", err)
	}
	return client, nil
}

func (s *Geth) BlockNumber(client *rpc.Client) (hexutil.Uint64, error) {
	var blockNumber hexutil.Uint64
	if err := client.Call(&blockNumber, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil

}
func (s *Geth) L2BlockNumber(l1Client bind.ContractBackend, l2OutputOracleAddr common.Address, l1BlockNumber *big.Int) (*big.Int, error) {
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
