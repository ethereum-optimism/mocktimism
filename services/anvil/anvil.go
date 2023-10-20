package anvil

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum-optimism/mocktimism/service-discovery"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type AnvilService struct {
	hostname    string
	port        int
	serviceType string
	id          string
	config      map[string]string
	cmd         *exec.Cmd
	logger      log.Logger
}

func NewAnvilService(logger log.Logger) *AnvilService {
	return &AnvilService{
		hostname:    "127.0.0.1",
		port:        8545,
		serviceType: "_anvil._tcp",
		id:          "anvilL1",
		config:      map[string]string{},
		logger:      logger,
	}
}

func (a *AnvilService) Hostname() string {
	return a.hostname
}

func (a *AnvilService) Port() int {
	return a.port
}

func (a *AnvilService) ServiceType() string {
	return a.serviceType
}

func (a *AnvilService) ID() string {
	return a.id
}

func (a *AnvilService) Config() servicediscovery.ServiceConfig {
	return a.config
}

func (a *AnvilService) Start(ctx context.Context) error {
	// TODO make sure this command exists in path https://github.com/ethereum-optimism/mocktimism/issues/61
	// TODO user should be able to configure where in path it is https://github.com/ethereum-optimism/mocktimism/issues/61
	a.cmd = exec.CommandContext(ctx, "anvil")

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
		return fmt.Errorf("failed to start Anvil: %v", err)
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
	client, err := rpc.Dial(fmt.Sprintf("http://%s:%d", a.hostname, a.port))
	if err != nil {
		return false, fmt.Errorf("failed to dial RPC: %v", err)
	}
	defer client.Close()
	var result string
	err = client.Call(&result, "eth_blockNumber")
	if err != nil {
		return false, fmt.Errorf("failed to retrieve block number: %v", err)
	}
	return result != "", nil
}
