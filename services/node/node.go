package node

import (
	"fmt"
	"os/exec"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
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

func (s *NodeService) Start() error {
	return nil
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

func (s *NodeService) CreateLifecycle(ctx *cli.Context) (cliapp.Lifecycle, error) {
	cfg, err := opnode.NewConfig(ctx, s.logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create the rollup node config: %w", err)
	}
	// cfg.Cancel = closeApp

	snapshotLog, err := opnode.NewSnapshotLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create snapshot root logger: %w", err)
	}

	n, err := node.New(
		ctx.Context,
		cfg,
		s.logger,
		snapshotLog,
		"TODO add version",
		metrics.NewMetrics("default"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create the rollup node: %w", err)
	}

	return n, nil
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
