package anvil

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum-optimism/mocktimism/config"
	"github.com/stretchr/testify/require"
)

func TestAnvilService_Run(t *testing.T) {
	// Create a temporary Chain configuration
	chainConfig := config.Chain{
		ForkURL:      "https://mainnet.optimism.io",
		Port:         8547, // use a different port for testing to avoid conflicts
		Host:         "127.0.0.1",
		BlockTime:    12,
		PruneHistory: false,
	}

	fmt.Println("creating anvil")
	// Initialize a new Anvil service
	anvilService := NewAnvil(chainConfig)

	// Run the service in a separate goroutine
	go func() {
		fmt.Println("Running anvil")
		anvilService.Run()
	}()

	// Wait for the service to fully start. In a real-world scenario, you might want to implement
	// a more robust mechanism for ensuring the service is ready (e.g., a health check).
	time.Sleep(10 * time.Second)

	// Issue a command to get the block_number
	cmd := exec.Command(
		"cast",
		"block_number",
		"http://"+chainConfig.Host+":"+strconv.FormatInt(int64(chainConfig.Port), 10),
	)
	output, err := cmd.Output()
	require.NoError(t, err)

	// Assert the expected block number
	expectedBlockNumber := "0"
	require.Equal(t, expectedBlockNumber, string(output))
}
