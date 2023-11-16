package wireguard

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/jcocozza/wg-monitor/utils"
)

// run wg show on specific interface
func WgSpecific(interfaceName string) []byte {
	cmd := exec.Command("wg", "show", interfaceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		slog.Error("Failed to run wg show "+interfaceName)
		panic(err)
	}

	return output
}

// Read in a wireguard config file
func readConf(configurationPath string) []byte {
	conf, err := utils.ReadFile(configurationPath)

	if err != nil {
		slog.Error("Failed to read configuration file for "+configurationPath)
		panic(err)
	}
	return conf
}

// wireguardPath -- path to folder;
// Get the .conf files in the wireguardPath
func getConfNames(wireguardPath string) []byte {
	cmd1 := exec.Command("ls", wireguardPath)
	cmd2 := exec.Command("grep", ".conf")

	// Get the output of the first command
	outCmd1, err := cmd1.Output()
	if err != nil {
		slog.Error("Error running ls command:", err)
		panic(err) 
	}
	// Set the input of the second command to the output of the first command
	cmd2.Stdin = strings.NewReader(string(outCmd1))

	// Get the output of the second command
	outCmd2, err := cmd2.Output()
	if err != nil {
		slog.Error("Error running grep command:", err)
		panic(err)
	}
	return outCmd2
}

// generate private key 
func wgGenKey() []byte {
	cmd := exec.Command("wg", "genkey")
	privateKey, err := cmd.Output()

	if err != nil {
		slog.Error("Failed to generate private key")
	}

	slog.Debug("Generated Public Key:"+string(privateKey))
	return privateKey
}

// generate public key from private key
func wgPubKey(privateKey string) []byte {
	pubKeyCmd := exec.Command("wg", "pubkey")
	pubKeyCmd.Stdin = strings.NewReader(string(privateKey)) // pass 
	publicKey, err := pubKeyCmd.Output()

	if err != nil {
		slog.Error("Failed to generate public key")
		panic(err)
	}

	slog.Debug("Generated Public Key:"+string(publicKey))
	return publicKey
}

// reload the server without a reboot (useful when adding new peers)
func ReloadServer(interfaceName string) {
	cmdString := fmt.Sprintf("sudo wg syncconf %s <(sudo wg-quick strip %s)", interfaceName, interfaceName) 
	cmd := exec.Command(cmdString)

	err := cmd.Run()

	if err != nil {
		slog.Error("Failed to reload server: "+interfaceName)
	}
}
