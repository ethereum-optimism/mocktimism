package mocktimism

import (
	"context"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
)

type Mocktimism struct {
	log log.Logger
	// shutdown requests the service that maintains the indexer to shut down,
	// and provides the error-cause of the critical failure (if any).
	shutdown context.CancelCauseFunc

	stopped atomic.Bool
}

func NewMocktimism(
	ctx context.Context,
	log log.Logger,
	shutdown context.CancelCauseFunc,
) (*Mocktimism, error) {
	return &Mocktimism{log: log, shutdown: shutdown}, nil
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
