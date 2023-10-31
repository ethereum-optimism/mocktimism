package anvil

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum-optimism/mocktimism/config"
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

func (a *AnvilService) Port() uint {
	return a.config.Port
}

func (a *AnvilService) ServiceType() string {
	return SERVICE_TYPE
}

func (a *AnvilService) ID() string {
	return a.id
}

func (a *AnvilService) Config() config.Chain {
	return a.config
}

func (a *AnvilService) Start(ctx context.Context) error {
	// TODO make sure this command exists in path https://github.com/ethereum-optimism/mocktimism/issues/61
	// TODO user should be able to configure where in path it is https://github.com/ethereum-optimism/mocktimism/issues/61
	a.cmd = exec.CommandContext(ctx, "anvil", "--port", fmt.Sprintf("%d", a.config.Port), "--host", a.config.Host)

	stdout, _ := a.cmd.StdoutPipe()
	stderr, _ := a.cmd.StderrPipe()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			a.logger.Info(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			a.logger.Error(scanner.Text())
		}
	}()

	err := a.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start Anvil: %w", err)
	}
	a.logger.Info("Started Anvil...")
	if err := a.cmd.Wait(); err != nil {
		a.logger.Error("Anvil process terminated with an error", "error", err)
	} else {
		a.logger.Info("Anvil process terminated normally")
	}
	return nil
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
	defer client.Close()

	var result string
	err = client.Call(&result, "eth_blockNumber")
	if err != nil {
		return false, fmt.Errorf("failed to retrieve block number: %w", err)
	}
	return result != "", nil
}
