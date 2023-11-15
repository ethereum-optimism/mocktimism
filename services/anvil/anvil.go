// TODO generalize most of this boilerplate into a common service interface
package anvil

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum-optimism/mocktimism/process"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	SERVICE_TYPE = "anvil"
)

type AnvilService struct {
	id     string
	config config.Chain
	cmd    *exec.Cmd
	logger log.Logger
}

func validateConfig(cfg config.Chain) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port is required")
	}
	return nil
}

func NewAnvilService(id string, logger log.Logger, cfg config.Chain) (*AnvilService, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &AnvilService{
		id:     id,
		config: cfg,
		logger: logger,
	}, nil
}

func (a *AnvilService) Hostname() string {
	return a.config.Host
}

func (a *AnvilService) Port() int {
	return int(a.config.Port)
}

func (a *AnvilService) ServiceType() string {
	return SERVICE_TYPE
}

func (a *AnvilService) ID() string {
	return a.id
}

func (a *AnvilService) Config() interface{} {
	return a.config
}

func (a *AnvilService) Start(ctx context.Context) error {
	a.cmd = exec.CommandContext(ctx, "anvil", configToArgs(a.config)...)
	return process.RunCommand(ctx, a.cmd, a.logger, "anvil")
}

func (a *AnvilService) Stop() error {
	if a.cmd == nil {
		return fmt.Errorf("Service is not running")
	}
	return a.cmd.Cancel()
}

func (a *AnvilService) HealthCheck() (bool, error) {
	client, err := rpc.Dial(fmt.Sprintf("http://%s:%d", a.config.Host, a.config.Port))
	if err != nil {
		return false, fmt.Errorf("failed to dial RPC: %w", err)
	}
	if err != nil {
		return false, err
	}
	defer client.Close()
	if err != nil {
		return false, err
	}
	_, err = blockNumber(client)
	if err != nil {
		return false, err
	}
	return true, nil
}
