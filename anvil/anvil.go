package anvil

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ethereum-optimism/mocktimism/config"
	"go-micro.dev/v4"
)

type Anvil struct {
	service micro.Service
	config  config.Chain
}

var (
	serviceName = "anvil"
	version     = "latest"
)

func NewAnvil(config config.Chain) *Anvil {
	// how do we deal with log levels?
	// how do we add a healthcheck?
	service := micro.NewService()
	fmt.Println("initing...")
	metadata := chainToMetadata(config)
	a := Anvil{
		service: service,
		config:  config,
	}
	// what is cleanest way to hook this up? Need to find an example
	a.service.Init(
		micro.Name(serviceName),
		micro.Address(fmt.Sprintf(":%d", config.Port)),
		micro.Metadata(metadata),
		micro.Version(version),
	)

	return &a
}

func chainToMetadata(chain config.Chain) map[string]string {
	return map[string]string{
		"type":     "anvil",
		"fork_url": chain.ForkURL,
		// TODO add more fileds here https://github.com/ethereum-optimism/mocktimism/issues/28
	}
}

// doesn't work yet
func (a *Anvil) Run() {
	cmdStr := constructAnvilCommand(a.config)
	fmt.Println("running... ", cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to run anvil:", err)
		return
	}
	a.service.Run()
}

func constructAnvilCommand(config config.Chain) string {
	// TODO path to anvil should be configurable https://github.com/ethereum-optimism/mocktimism/issues/28
	// TODO should handle anvil not existing via throwing a useful error with link to docs https://github.com/ethereum-optimism/mocktimism/issues/28
	cmdParts := []string{"anvil"}
	if config.ForkURL != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("--fork-url %s", config.ForkURL))
	}
	fmt.Println(strings.Join(cmdParts, " "))
	// TODO Add other flags here as needed https://github.com/ethereum-optimism/mocktimism/issues/28
	return strings.Join(cmdParts, " ")
}
