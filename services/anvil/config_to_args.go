package anvil

import (
	"fmt"
	"github.com/ethereum-optimism/mocktimism/config"
)

func configToArgs(c config.Chain) []string {
	args := []string{}

	if c.Port != 0 {
		args = append(args, "--port", fmt.Sprintf("%d", c.Port))
	}
	if c.Host != "" {
		args = append(args, "--host", c.Host)
	}
	if c.ForkBlockNumber != 0 {
		args = append(args, "--fork-block-number", fmt.Sprintf("%d", c.ForkBlockNumber))
	}
	if c.ForkChainID != 0 {
		args = append(args, "--fork-chain-id", fmt.Sprintf("%d", c.ForkChainID))
	}
	if c.ForkURL != "" {
		args = append(args, "--fork-url", c.ForkURL)
	}
	return args
}
