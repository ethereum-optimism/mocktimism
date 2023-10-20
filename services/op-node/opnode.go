package opnode

import (
	"context"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

type OpNodeService struct {
	hostname    string
	port        int
	serviceType string
	id          string
	config      map[string]string
	client      *http.Client
	node        *node.Node
	logger      log.Logger
}

func NewOpNodeService(logger log.Logger) *OpNodeService {
	return &OpNodeService{
		hostname:    "127.0.0.1",
		port:        9000,
		serviceType: "_opnode._tcp",
		id:          "opnode1",
		client:      &http.Client{},
		config:      map[string]string{},
		logger:      logger,
	}
}

func (o *OpNodeService) Start(ctx context.Context) error {
	// Here, we can call the `RollupNodeMain` function directly to start the node.
	// You might need to adjust this if you have specific parameters or context configurations.

	// We'll use a wrapped function to handle the initialization and starting of the op-node.
	rollupNodeInitializer := RollupNodeMain("version-info") // Adjust "version-info" with the version you want.
	err := rollupNodeInitializer(&cli.Context{})            // You may need to adjust this to pass a properly initialized context.

	return err
}

func (o *OpNodeService) HealthCheck() (bool, error) {
	resp, err := o.client.Get(fmt.Sprintf("http://%s:%d/healthz", o.hostname, o.port))
	if err != nil {
		return false, fmt.Errorf("failed to check health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("health check returned %d status", resp.StatusCode)
	}

	return true, nil
}
