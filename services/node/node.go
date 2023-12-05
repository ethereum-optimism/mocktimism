package node

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum-optimism/mocktimism/process"
	"github.com/ethereum/go-ethereum/log"
)

var (
	SERVICE_TYPE = "op-node"
	DEPENDS_ON   = []string{
		"l1",
		"l2",
	}
)

type NodeService struct {
	id     string
	config NodeConfig
	cmd    *exec.Cmd
	logger log.Logger
}

func NewNodeService(id string, logger log.Logger, cfg NodeConfig) (*NodeService, error) {
	return &NodeService{
		id:     id,
		config: cfg,
		logger: logger,
	}, nil
}

func (s *NodeService) Hostname() string {
	return s.config.Host
}

func (s *NodeService) Port() int {
	return int(s.config.Ports.RPCPort)
}

func (s *NodeService) ServiceType() string {
	return SERVICE_TYPE
}

func (s *NodeService) ID() string {
	return s.id
}

func (s *NodeService) Config() interface{} {
	return s.config
}

func (s *NodeService) HealthCheck() (bool, error) {
	// TODO hit healthz for healthcheck
	return true, nil
}

func (s *NodeService) Start(ctx context.Context) error {
	args := buildCommandArgs(s.config)
	s.cmd = exec.CommandContext(ctx, "op-node", args...)
	return process.RunCommand(ctx, s.cmd, s.logger, "op-node")
}

func (s *NodeService) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop Node service: %w", err)
		}
		s.logger.Info("message", "Node service stopped", "id", s.id)
	}
	return nil
}
