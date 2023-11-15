package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/ethereum/go-ethereum/log"
)

func RunCommand(
	ctx context.Context,
	cmd *exec.Cmd,
	logger log.Logger,
	name string,
) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	go logOutput(stdout, logger.Info)
	go logOutput(stderr, logger.Error)

	logger.Info("Starting service...", "service", name)
	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start service", "service", name, "error", err)
		return fmt.Errorf("failed to start service %s: %w", name, err)
	}
	logger.Info("Service started", "service", name)
	if err != nil {
		return fmt.Errorf("failed to start anvil: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		logger.Error("Anvil process terminated with an error", "error", err)
	} else {
		logger.Info("Anvil process terminated normally")
	}
	return nil
}

func logOutput(reader io.Reader, logMessage func(msg string, ctx ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logMessage(scanner.Text())
	}
}
