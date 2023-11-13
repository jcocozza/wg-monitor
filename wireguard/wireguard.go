package wireguard

import (
	"os/exec"
	"regexp"
	"strings"
)

//Show all interfaces and their peers
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

//parse output of wg()
func parseWgOut(wgInput []byte) []Interface {
	interfaces := []Interface{}
	//var interfaceData Interface
	wgOutput := string(wgInput)

	lines := strings.Split(wgOutput, "\n")
	for i, line := range lines {
		//slog.Debug(fmt.Sprintf("Parsing %d iteration", i))
		line = strings.TrimSpace(line)

		var name string
		var interfacePublicKey string
		var listeningPort string

		if strings.HasPrefix(line, "interface: ") {
			//because the interfaces and their peers are ordered, we can assume that when we get a new interface, we need to start over
			name = strings.TrimPrefix(line, "interface: ") //get name of interface

			interfacePublicKey = strings.TrimPrefix(lines[i+1],"public key: ")
			listeningPort = strings.TrimPrefix(lines[i+3], "listening port: ")	
			
			interfaceData := NewInterface(name, interfacePublicKey, listeningPort)
			interfaces = append(interfaces, *interfaceData)
			
		} else if strings.HasPrefix(line, "peer: ") {
			var endPoint string
			var ipList []string
			var latestHandshake string
			var transfer []string
			publicKey := strings.TrimPrefix(line, "peer: ")

			if strings.HasPrefix(strings.TrimSpace(lines[i+1]), "endpoint: ") {
				endPoint = strings.TrimSpace(lines[i+1])

				ips := strings.TrimPrefix(strings.TrimSpace(lines[i+2]), "allowed ips: ")
				ipList = strings.SplitAfter(ips, ",")

				latestHandshake = strings.TrimSpace(lines[i+3])

				transfers := strings.TrimSpace(lines[i+4])
				transfer = strings.SplitAfter(transfers, ",")


			} else { // if the peer has NOT connected to the server since it started
				endPoint = ""

				ips := strings.TrimPrefix(strings.TrimSpace(lines[i+1]), "allowed ips: ")
				ipList = strings.SplitAfter(ips, ",")

				latestHandshake = ""

				transfer = []string{}
			}

			// remove any extra ip stuff, so we are left with raw ip "10.5.5.1/32 -> 10.5.5.1"
			var allowedIPs []string
			for _, ipaddr := range ipList {
				cleanedIp := strings.Split(ipaddr, "/")[0]
				allowedIPs = append(allowedIPs, cleanedIp)	
			}

			peer := NewPeer(publicKey, endPoint, allowedIPs, latestHandshake, transfer)

			interfaces[len(interfaces) - 1].AddPeer(*peer)
		}
	}
	return interfaces
}

//wrapper for wg() and parseWgOut()
func GetWgInfo() ([]Interface) {
	wgOut, err := wg()

	if err != nil {
		panic(err)
	}
	interfaces := parseWgOut(wgOut)
	return interfaces
}

// simplier version to get just interface names
func GetInterfaceNames() []string {
	var interfaceNames []string
	wgOut, err := wg()

	if err != nil {
		panic(err)
	}
	info := string(wgOut)

	lines := strings.Split(info, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "interface: ") {
			interfaceNames = append(interfaceNames, strings.TrimPrefix(line, "interface: "))
		}
	}
	return interfaceNames
}

// Get an interface
func GetInterfaceByName(interfaceName string) Interface {
	wgOut, err := wgSpecific(interfaceName)

	if err != nil {
		panic(err)
	}

	interfaces := parseWgOut(wgOut)
	return interfaces[0] // there will only be 1 interface
}

// the the interface IP -- this will be what is set in the wireguard conf file
func GetInterfaceIP(interfaceName string) string {
	cmd := exec.Command("ifconfig", interfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	ifconfigResult := string(output)

	// Regular expression to match an IPv4 address
	ipRegex := regexp.MustCompile(`(\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)`)

	// Find all matches in the text
	ipMatches := ipRegex.FindAllString(ifconfigResult, -1)

	// Display the matched IP address
	if len(ipMatches) > 0 {
		ipAddress := ipMatches[0] // Assuming the first match is the IP address
		return ipAddress
	} else {
		panic("No IP Found -- is your wireguard server down?")
	}
}