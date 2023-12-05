package geth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
)

func initGethDataDir(genesis core.Genesis) error {
	// Create a new node configuration
	nodeConfig := node.DefaultConfig
	nodeConfig.DataDir = "/data" // The path to the data directory

	// Create a new node instance
	n, err := node.New(&nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create new node: %w", err)
	}

	// Register the Ethereum service
	ethConfig := eth.DefaultConfig
	ethConfig.Genesis = &genesis
	if err := n.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &ethConfig)
	}); err != nil {
		return fmt.Errorf("failed to register Ethereum service: %w", err)
	}

	// Initialize the Geth database with the genesis block
	db := n.OpenDatabase("chaindata", 0, 0, "")
	defer db.Close()
	genesisBlock, _ := genesis.ToBlock(db)
	if genesisBlock != nil {
		core.WriteHeadBlockNumber(db, genesisBlock.NumberU64())
		core.WriteCanonicalHash(db, genesisBlock.Hash(), genesisBlock.NumberU64())
		core.WriteHeader(db, genesisBlock.Header())
	}

	return nil
}
