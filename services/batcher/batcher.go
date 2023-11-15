package batcher

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum/go-ethereum/log"
)

type OpBatcherService struct {
	id     string
	config OpBatcherConfig
	cmd    *exec.Cmd
	logger log.Logger
}

func NewOpBatcherService(id string, logger log.Logger, cfg OpBatcherConfig) (*OpBatcherService, error) {
	return &OpBatcherService{
		id:     id,
		config: cfg,
		logger: logger,
	}, nil
}

func (s *OpBatcherService) Start(ctx context.Context) error {
	args := s.buildCommandArgs()
	s.cmd = exec.CommandContext(ctx, "op-batcher", args...)

	// Handle STDOUT and STDERR similarly to other services

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OpBatcher service: %w", err)
	}

	s.logger.Info("OpBatcher service started", "id", s.id)
	return nil
}

func (s *OpBatcherService) Stop() error {
	// Implement stop logic similar to other services
}

func (s *OpBatcherService) buildCommandArgs() []string {
	// Construct command arguments from the configuration
	// e.g., "--l1-eth-rpc=" + s.config.L1EthRPC, etc.
}

// Implement other methods as needed
