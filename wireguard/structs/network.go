package structs

import (
	"strconv"
	"strings"

	c "github.com/jcocozza/wg-monitor/wireguard/commands"
)

type NetworkInterface struct {
	Name 			string `json:"name"`
	PublicKey 		string `json:"publicKey"`
	ListeningPort 	int	   `json:"listeningPort"`
	Status			bool   `json:"status"`
}

func NewNetworkInterface(name string, publicKey string, listeningPort int, status bool) *NetworkInterface {
	return &NetworkInterface{
		Name: name,
		PublicKey: publicKey,
		ListeningPort: listeningPort,
		Status: status,
	}
}

func EmptyNetworkInterface() *NetworkInterface {
	return &NetworkInterface{
		Name: "NoNetwork",
		PublicKey: "",
		ListeningPort: -999,
		Status: false,
	}
}

// Check whether a given Network interface is running or not. If so return true.
func (network *NetworkInterface) CheckStatus() bool {
	result := c.WgSpecific(network.Name)

	if result != nil {
		network.Status = true
	} else {
		network.Status = false
	}
	return network.Status
}

// get all the aviable interfaces from `wg show`. 
// returns a map of publicKey -> NetworkInterface object
func LoadAllNetworkInterfaceInfo() map[string]*NetworkInterface {
	networkInfoString := string(c.WgShow())

	networkMap := make(map[string]*NetworkInterface)

	var networkInterface *NetworkInterface
	var interfaceName string
	var publicKey string
	var listeningPort int
	
	lines := strings.Split(networkInfoString, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "interface:") {
			interfaceName = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "interface: "))
			publicKey     = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i+1]), "public key: "))
			listeningPort, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i+3]), "listening port: ")))
			
			networkInterface = NewNetworkInterface(interfaceName, publicKey, listeningPort, true)
			networkMap[publicKey] = networkInterface
		}

	}
	return networkMap
}