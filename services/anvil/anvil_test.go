package anvil

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestAnvilServiceValidation(t *testing.T) {
	invalidCfgs := []config.Chain{
		// No PORT
		{
			Host: "127.0.0.1",
		},
		// No HOST
		{
			Port: 8545,
		},
	}
	for _, cfg := range invalidCfgs {
		_, err := NewAnvilService("TestService", log.New("module", "test"), cfg)
		require.Error(t, err)
	}
}

func TestAnvilService(t *testing.T) {
	logger := log.New("module", "test")
	cfg := config.Chain{
		Host: "127.0.0.1",
		Port: 8545,
	}

	// Initialize the AnvilService
	service, err := NewAnvilService("TestService", logger, cfg)
	require.NoError(t, err, "Failed to initialize the Anvil service")

	// Start the service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // timeout to ensure the service doesn't run indefinitely
	defer cancel()
	err = service.Start(ctx)
	require.NoError(t, err, "Failed to start the Anvil service")

	// Poll for health check until healthy or timeout
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond) // polling every 200ms
	defer ticker.Stop()

	healthy := false
loop:
	for {
		select {
		case <-timeout:
			break loop
		case <-ticker.C:
			healthy, err = service.HealthCheck()
			if healthy && err != nil {
				break loop
			}
		}
	}

	require.NoError(t, err, "Health check failed")
	require.True(t, healthy, "Service is not healthy after waiting for 2 seconds")

	// Verify service details
	require.Equal(t, "127.0.0.1", service.Hostname())
	require.Equal(t, uint(8545), service.Port())
	require.Equal(t, "anvil", service.ServiceType())

	// Stop the service
	err = service.Stop()
	require.NoError(t, err, "Failed to stop the Anvil service")
}

func TestStopWithoutStarting(t *testing.T) {
	logger := log.New("module", "test")
	cfg := config.Chain{
		Host: "127.0.0.1",
		Port: 8545,
	}
	service, err := NewAnvilService("TestService", logger, cfg)
	require.NoError(t, err, "Failed to initialize the Anvil service")

	err = service.Stop()
	require.Error(t, err, "Expected an error when stopping a service that hasn't been started")
}
