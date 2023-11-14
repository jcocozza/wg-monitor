package wireguard

import (
	"os/exec"
	"strings"
)

// show all interfaces and their peers
func wg() ([]byte, error ){
	cmd := exec.Command("wg", "show", "all")
	output, err := cmd.CombinedOutput()

	return output, err
}

// run wg show on specific interface
func wgSpecific(interfaceName string) ([]byte, error) {
	cmd := exec.Command("wg", "show", interfaceName)
	output, err := cmd.CombinedOutput()

	return output, err
}

// generate private key 
func wgGenKey() ([]byte, error){
	cmd := exec.Command("wg", "genkey")
	privateKey, err := cmd.Output()
	return privateKey, err
}

// generate public key from private key
func wgPubKey(privateKey string) ([]byte, error) {
	pubKeyCmd := exec.Command("wg", "pubkey")
	pubKeyCmd.Stdin = strings.NewReader(string(privateKey)) // pass 
	publicKey, err := pubKeyCmd.Output()
	return publicKey, err
}

// run ifconfig on a given interface
func ifconfig(interfaceName string) ([]byte, error){
	cmd := exec.Command("ifconfig", interfaceName)
	output, err := cmd.CombinedOutput()

	return output, err
}