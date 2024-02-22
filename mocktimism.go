package mocktimism

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/services/anvil"
	"github.com/ethereum-optimism/mocktimism/services/batcher"
	"github.com/ethereum-optimism/mocktimism/services/node"
	"github.com/ethereum-optimism/mocktimism/services/proposer"
	"github.com/ethereum/go-ethereum/log"
)

type Mocktimism struct {
	log log.Logger
	// shutdown requests the service that maintains the indexer to shut down,
	// and provides the error-cause of the critical failure (if any).
	shutdown context.CancelCauseFunc

	cfg *config.Config

	stopped atomic.Bool

	// services
	anvilL1    *anvil.AnvilService
	anvilL2    *anvil.AnvilService
	opNode     *node.NodeService
	opProposer *proposer.ProposerService
	opBatcher  *batcher.BatcherService
}

func NewMocktimism(
	ctx context.Context,
	log log.Logger,
	shutdown context.CancelCauseFunc,
	cfg *config.Config,
) (*Mocktimism, error) {
	out := &Mocktimism{log: log, shutdown: shutdown, cfg: cfg}
	if err := out.initAnvilL1(ctx); err != nil {
		return nil, err
	}
	if err := out.initAnvilL2(ctx); err != nil {
		return nil, err
	}
	if err := out.initOpNode(ctx); err != nil {
		return nil, err
	}
	if err := out.initOpProposer(ctx); err != nil {
		return nil, err
	}
	if err := out.initOpBatcher(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *Mocktimism) Start(ctx context.Context) error {
	m.log.Debug("starting anvil-l1...")
	if err := m.anvilL1.Start(ctx); err != nil {
		return fmt.Errorf("failed to start anvil-l1: %w", err)
	}
	m.log.Debug("starting anvil-l2...")
	if err := m.anvilL2.Start(ctx); err != nil {
		return fmt.Errorf("failed to start anvil-l2: %w", err)
	}
	m.log.Debug("starting op-node...")
	if err := m.opNode.Start(); err != nil {
		return fmt.Errorf("failed to start op-node: %w", err)
	}
	m.log.Debug("starting op-proposer...")
	if err := m.opProposer.Start(); err != nil {
		return fmt.Errorf("failed to start op-proposer: %w", err)
	}
	m.log.Debug("starting op-batcher...")
	if err := m.opBatcher.Start(); err != nil {
		return fmt.Errorf("failed to start op-batcher: %w", err)
	}
	m.log.Debug("mocktimism started")
	return nil
}

func (m *Mocktimism) Stop(ctx context.Context) error {
	var result error

	m.log.Debug("stopping mocktimism...")

	if m.anvilL1 != nil {
		if err := m.anvilL1.Stop(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop anvil-l1: %w", err))
		}
	}
	if m.anvilL2 != nil {
		if err := m.anvilL2.Stop(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop anvil-l2: %w", err))
		}
	}
	if m.opNode != nil {
		if err := m.opNode.Stop(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop op-node: %w", err))
		}
	}
	if m.opProposer != nil {
		if err := m.opProposer.Stop(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop op-proposer: %w", err))
		}
	}
	if m.opBatcher != nil {
		if err := m.opBatcher.Stop(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop op-batcher: %w", err))
		}
	}

	m.stopped.Store(true)

	if result == nil {
		m.log.Debug("mocktimism stopped successfully")
	} else {
		m.log.Error("failed to stop mocktimism", "errors", result)
	}

	return result
}

func (m *Mocktimism) Stopped() bool {
	return m.stopped.Load()
}

func (m *Mocktimism) getL1Chain(ctx context.Context) (config.Chain, error) {
	// we are assuming only 1 profile 1ith only 1 l1 chain as of now
	for _, profile := range m.cfg.Profiles {
		for _, chain := range profile.Chains {
			if chain.BaseChainID == chain.ChainID {
				return chain, nil
			}
		}
	}
	return config.Chain{}, errors.New("no l1 chain found")
}
func (m *Mocktimism) getL2Chain(ctx context.Context) (config.Chain, error) {
	// we are assuming only 1 profile 1ith only 1 l2 chain as of now
	for _, profile := range m.cfg.Profiles {
		for _, chain := range profile.Chains {
			if chain.BaseChainID != chain.ChainID {
				return chain, nil
			}
		}
	}
	return config.Chain{}, errors.New("no l1 chain found")
}

func (m *Mocktimism) initAnvilL1(ctx context.Context) error {
	chain, err := m.getL1Chain(ctx)
	if err != nil {
		return err
	}
	s, err := anvil.NewAnvilService("anvil-l1", m.log, chain)
	if err != nil {
		return err
	}
	m.anvilL1 = s
	return nil
}

func (m *Mocktimism) initAnvilL2(ctx context.Context) error {
	chain, err := m.getL2Chain(ctx)
	if err != nil {
		return err
	}
	s, err := anvil.NewAnvilService("anvil-l2", m.log, chain)
	if err != nil {
		return err
	}
	m.anvilL2 = s
	return nil
}

func (m *Mocktimism) initOpNode(ctx context.Context) error {
	s, err := node.NewNodeService("op-node", m.log, node.NodeConfig{})
	if err != nil {
		return err
	}
	m.opNode = s
	return nil
}

func (m *Mocktimism) initOpProposer(ctx context.Context) error {
	s, err := proposer.NewProposerService("proposer", m.log, proposer.ProposerConfig{})
	if err != nil {
		return err
	}
	m.opProposer = s
	return nil
}

func (m *Mocktimism) initOpBatcher(ctx context.Context) error {
	s, err := batcher.NewBatcherService("batcher", m.log, batcher.BatcherConfig{})
	if err != nil {
		return err
	}
	m.opBatcher = s
	return nil
}
