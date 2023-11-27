package commands

import (
	"log/slog"
	"os/exec"
	"strings"
)

// wireguardPath -- path to folder;
// Get the .conf files in the wireguardPath
func GetConfNames(wireguardPath string) []byte {
	slog.Debug("getting configuration names...")
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

// generate a private, public key pairing
// the basic idea is: "wg genkey | tee client_privatekey | wg pubkey > client_publickey"
func GenerateKeyPair() (string, string){
	slog.Debug("generating keypair...")
	privateKey := WgGenKey()
	publicKey := WgPubKey(string(privateKey))

	return strings.TrimSpace(string(privateKey)), strings.TrimSpace(string(publicKey))
}

