package proposer

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
	SERVICE_TYPE = "op-proposer"
	DEPENDS_ON   = []string{
		"l1",
		"l2",
	}
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

func (s *ProposerService) Start() error {
	return nil
}

func (s *ProposerService) Hostname() string {
	return s.config.Host
}

func (s *ProposerService) Port() int {
	return int(s.config.Ports.RPCPort)
}

func (s *ProposerService) ServiceType() string {
	return SERVICE_TYPE
}

func (s *ProposerService) ID() string {
	return s.id
}

func (s *ProposerService) Config() interface{} {
	return s.config
}

func (s *ProposerService) HealthCheck() (bool, error) {
	// TODO hit healthz for healthcheck
	return true, nil
}

func (s *ProposerService) CreateLifecycle(ctx *cli.Context) (cliapp.Lifecycle, error) {
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

func (s *ProposerService) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop Proposer service: %w", err)
		}
		s.logger.Info("message", "Proposer service stopped", "id", s.id)
	}
	return nil
}
