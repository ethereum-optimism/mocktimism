package prysmmanager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/ethereum/go-ethereum/core"
)

const prysmVersion = "v4.0.8"

//go:embed static/prysm/config.yml
var prysmConfigFile []byte

func downloadBin(binaryName string) error {
	arch := runtime.GOARCH
	url := fmt.Sprintf("https://github.com/prysmaticlabs/prysm/releases/download/%s/%s-%s-linux-%s", prysmVersion, binaryName, prysmVersion, arch)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("/usr/local/bin/" + binaryName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func savePrysmConfig() (string, error) {
	tmpFile, err := os.CreateTemp("", "prysm-config-*.yml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(prysmConfigFile)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func InstallPrysm() error {
	binaries := []string{"prysmctl", "validator", "beacon-chain"}
	for _, bin := range binaries {
		err := downloadBin(bin)
		if err != nil {
			return err
		}
		err = os.Chmod("/usr/local/bin/"+bin, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// You can add this in the same file or a separate one within the same package

func startPrysmComponent(component string, args ...string) error {
	cmd := exec.Command("/usr/local/bin/"+component, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func saveGethConfig(genesis core.Genesis) (string, error) {
	configBytes, err := json.Marshal(genesis)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "geth-config-*.json")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(configBytes)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func StartBeaconChain(genesis core.Genesis) error {
	prysmConfigPath, err := savePrysmConfig()
	if err != nil {
		return err
	}

	gethConfigPath, err := saveGethConfig(genesis)

	if err != nil {
		return err
	}

	args := []string{
		"--config-file=" + prysmConfigPath,
		"--geth-config-file=" + gethConfigPath, // Assuming Prysm accepts Geth config this way
	}
	return startPrysmComponent("beacon-chain", args...)
}

func StartValidator(args ...string) error {
	return startPrysmComponent("validator", args...)
}
