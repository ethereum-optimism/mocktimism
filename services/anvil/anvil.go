package anvil

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	args := []string{}

	if a.config.Port != 0 {
		args = append(args, "--port", fmt.Sprintf("%d", a.config.Port))
	}
	if a.config.Host != "" {
		args = append(args, "--host", a.config.Host)
	}
	if a.config.ForkBlockNumber != 0 {
		args = append(args, "--fork-block-number", fmt.Sprintf("%d", a.config.ForkBlockNumber))
	}
	if a.config.ForkChainID != 0 {
		args = append(args, "--fork-chain-id", fmt.Sprintf("%d", a.config.ForkChainID))
	}
	if a.config.ForkURL != "" {
		args = append(args, "--fork-url", a.config.ForkURL)
	}

	a.cmd = exec.CommandContext(ctx, "anvil", args...)

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

	a.logger.Info("Starting Anvil...")
	err := a.cmd.Start()
	if err != nil {
		a.logger.Error("Failed to start anvil: ", err)
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
	client, err := a.GetClient()
	if err != nil {
		return false, err
	}
	defer client.Close()
	if err != nil {
		return false, err
	}
	_, err = a.BlockNumber(client)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *AnvilService) GetClient() (*rpc.Client, error) {
	client, err := rpc.Dial(fmt.Sprintf("http://%s:%d", a.config.Host, a.config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC: %w", err)
	}
	return client, nil
}

func (a *AnvilService) BlockNumber(client *rpc.Client) (hexutil.Uint64, error) {
	var blockNumber hexutil.Uint64
	if err := client.Call(&blockNumber, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}
