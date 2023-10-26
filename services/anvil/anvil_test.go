package anvil

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestAnvilService(t *testing.T) {
	logger := log.New("module", "test")

	// Initialize the AnvilService
	service := NewAnvilService(logger)

	// Start the service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // timeout to ensure the service doesn't run indefinitely
	defer cancel()
	err := service.Start(ctx)
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
			if healthy || err != nil {
				break loop
			}
		}
	}

	require.NoError(t, err, "Health check failed")
	require.True(t, healthy, "Service is not healthy after waiting for 2 seconds")

	// Verify service details
	require.Equal(t, "127.0.0.1", service.Hostname())
	require.Equal(t, 8545, service.Port())
	require.Equal(t, "_anvil._tcp", service.ServiceType())

	// Stop the service
	err = service.Stop()
	require.NoError(t, err, "Failed to stop the Anvil service")
}

func TestStopWithoutStarting(t *testing.T) {
	logger := log.New("module", "test")
	service := NewAnvilService(logger)

	err := service.Stop()
	require.Error(t, err, "Expected an error when stopping a service that hasn't been started")
}
