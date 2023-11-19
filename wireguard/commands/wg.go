package commands

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// run wg show on specific interface
func WgSpecific(interfaceName string) []byte {
	cmd := exec.Command("wg", "show", interfaceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// might not be an error, if the interface isn't up, there will be an error.
		slog.Debug("Failed to run wg show "+ interfaceName)
		return nil
	}

	return output
}


func WgShow() []byte {
	cmd := exec.Command("wg", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to run wg show")
		return nil
	}
	return output
}


// generate private key 
func WgGenKey() []byte {
	cmd := exec.Command("wg", "genkey")
	privateKey, err := cmd.Output()

	if err != nil {
		slog.Error("Failed to generate private key")
	}

	slog.Debug("Generated Public Key:"+string(privateKey))
	return privateKey
}

// generate public key from private key
func WgPubKey(privateKey string) []byte {
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
func WgReloadServer(confName string) {
	cmdString := fmt.Sprintf("sudo wg syncconf %s <(sudo wg-quick strip %s)", confName, confName) 
	cmd := exec.Command(cmdString)

	err := cmd.Run()

	if err != nil {
		slog.Error("Failed to reload server: "+confName)
	}
}

