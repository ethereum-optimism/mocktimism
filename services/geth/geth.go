package geth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

type GethConfig struct {
	DataDir   string
	Verbosity int
	HTTPPort  int
	OpGeth    bool
}

type Geth struct {
	log    log.Logger
	config GethConfig

	node *node.Node
	eth  *eth.Ethereum
}

var (
	blockSignerAddress = common.Address{}
	gasCeilL1          = uint64(8000000)
	testJwtSecret      = "688f5d737bad920bdfb2fc2f488d6b6209eebda1dae949a8de91398d932c517a"
)

func NewGeth(name string, logger log.Logger, cfg GethConfig, genesis *core.Genesis) (*Geth, error) {
	// TODO we can delete handling genesis being nil once genesis is implemented https://github.com/ethereum-optimism/mocktimism/issues/86
	if genesis == nil {
		genesis = core.DefaultGenesisBlock()
	}

	var ethCfg *ethconfig.Config
	var formattedName string
	if cfg.OpGeth {
		formattedName = fmt.Sprintf("l2-geth-%v", name)
		ethCfg = &ethconfig.Config{
			// TODO some of these should be read from the toml config
			NetworkId: genesis.Config.ChainID.Uint64(),
			Genesis:   genesis,
			SyncMode:  downloader.FullSync,
			// Warning archive node is required or else trie nodes are pruned within minutes of starting devnet
			NoPruning: true,
			Miner: miner.Config{
				// l2 shouldn't mine
				Etherbase:         common.Address{},
				ExtraData:         nil,
				Recommit:          0,
				NewPayloadTimeout: 0,
				GasCeil:           0,
				GasFloor:          0,
				GasPrice:          nil,
			},
		}
	} else {
		formattedName = fmt.Sprintf("l1-geth-%v", name)
		ethCfg = &ethconfig.Config{
			// TODO some of these should be read from the toml config
			NetworkId: genesis.Config.ChainID.Uint64(),
			Genesis:   genesis,
			SyncMode:  downloader.FullSync,
			// Warning archive node is required or else trie nodes are pruned within minutes of starting devnet
			NoPruning: true,
			Miner: miner.Config{
				Etherbase:         blockSignerAddress,
				ExtraData:         nil,
				Recommit:          0,
				NewPayloadTimeout: 0,
				GasCeil:           gasCeilL1,
				GasFloor:          0,
				GasPrice:          big.NewInt(1),
			},
		}
	}

	nodeCfg := &node.Config{
		Name: formattedName,
		// TODO hardcoded these need to read the chain config instead
		HTTPHost:    "127.0.0.1",
		HTTPPort:    cfg.HTTPPort,
		WSHost:      "127.0.0.1",
		WSPort:      0,
		AuthAddr:    "127.0.0.1",
		AuthPort:    0,
		HTTPModules: []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
		WSModules:   []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
		// TODO I think we don't need this but leaving this here for now for reference. Delete if not needed when cleaning up
		// WSPort: cfg.WSPort,
		DataDir: cfg.DataDir,
		P2P: p2p.Config{
			NoDiscovery: true,
			// For OP devnet max peers for l2 is 0 and 1 for l1. I don't think this matters though
			MaxPeers: 1,
		},
		JWTSecret: testJwtSecret,
	}

	n, err := node.New(nodeCfg)

	// TODO e2e utils call n.Merger().FinalizePos(). I don't think we need this. Delete this comment if not needed.
	// e2e utils also run a fakePos via l1Node.RegisterLifecycle. This I also do not believe we need.
	if err != nil {
		return nil, err
	}

	backend, err := eth.New(n, ethCfg)
	if err != nil {
		return nil, err
	}

	// e2e utils call catalyst.Register(l2Node, backend) on the l2 node to enable engine api.
	// I don't think we need this because we have engine enabled in HTTPModules but leaving this note here just in case
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
	if s.config.OpGeth {
		return nil, fmt.Errorf("L2BlockNumber is meant to be called on l1 chains")
	}
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
