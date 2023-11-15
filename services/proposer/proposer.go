package proposer

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum/go-ethereum/log"
)

type ProposerService struct {
	id     string
	config ProposerConfig
	cmd    *exec.Cmd
	logger log.Logger
}

func NewProposerService(id string, logger log.Logger, cfg ProposerConfig) (*ProposerService, error) {
	return &ProposerService{
		id:     id,
		config: cfg,
		logger: logger,
	}, nil
}

func (s *ProposerService) Start(ctx context.Context) error {
	args := s.buildCommandArgs()
	s.cmd = exec.CommandContext(ctx, "op-proposer", args...)

	// Handle STDOUT and STDERR similar to NodeService

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Proposer service: %w", err)
	}

	s.logger.Info("Proposer service started", "id", s.id)
	return nil
}

func (s *ProposerService) Stop() error {
	// Implement stop logic similar to NodeService
}

func (s *ProposerService) buildCommandArgs() []string {
	// Construct command arguments from the configuration
	// e.g., "--poll-interval=" + s.config.PollInterval, etc.
}

// Implement other methods as needed
