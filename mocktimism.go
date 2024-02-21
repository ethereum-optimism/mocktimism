package mocktimism

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/services/anvil"
	"github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer"
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
	opNode     *node.OpNode
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
	if err := out.initServices(ctx); err != nil {
		return nil, errors.Join(err, out.Stop(ctx))
	}
	return out, nil
}

func (m *Mocktimism) Start(ctx context.Context) error {
	return nil
}

func (m *Mocktimism) Stop(ctx context.Context) error {
	if m.stopped.Load() {
		return nil
	}
	m.stopped.Store(true)
	return nil
}

func (m *Mocktimism) Stopped() bool {
	return m.stopped.Load()
}

// Inits every service in mocktimism
// we don't worry about initing the challenger because we are not using it
func (m *Mocktimism) initServices(ctx context.Context) error {
	if err := m.initAnvilL1(ctx); err != nil {
		return err
	}
	if err := m.initAnvilL2(ctx); err != nil {
		return err
	}
	if err := m.initOpNode(ctx); err != nil {
		return err
	}
	if err := m.initOpProposer(ctx); err != nil {
		return err
	}
	if err := m.initOpBatcher(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Mocktimism) initAnvilL1(ctx context.Context) error {
	return nil
}

func (m *Mocktimism) initAnvilL2(ctx context.Context) error {
	return nil
}

func (m *Mocktimism) initOpNode(ctx context.Context) error {
	return nil
}

func (m *Mocktimism) initOpProposer(ctx context.Context) error {
	return nil
}

func (m *Mocktimism) initOpBatcher(ctx context.Context) error {
	return nil
}
