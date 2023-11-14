package wireguard

import (
	"log/slog"
	"strconv"
	"strings"
)

// Check all the .conf files in the wireguard path.
// Create configurations for each of them
func InitWireGuardConfigurations(wireguardPath string) *WireGuardConfigurations {
	files := strings.TrimSpace(string(getConfNames(wireguardPath)))
	fileList := strings.Split(files, "\n")

	var configurationList []*Configuration
	for _, fileName := range fileList {
		slog.Debug("parsing conf:"+fileName)
		conf := parseConf(wireguardPath,fileName)
		configurationList = append(configurationList, conf)
	}

	wgcs := NewWireGuardConfigurations(configurationList)

	for _, confName := range wgcs.ConfNames {
		LoadPeerInfo(confName, wgcs)
	}

	return wgcs
}

// Parse a wireguard configuration file
func parseConf(configurationPath string, configurationFileName string) *Configuration {

	confByte := readConf(configurationPath+configurationFileName)
	confstr := string(confByte)

	confElements := strings.Split(confstr, "[Peer]") //split on "[Peer]"

	name, ifaceAddress, ifaceListenPort, ifacePrivateKey, ifaceDNS := parseInterfaceInfo(confElements[0]) // the first thing in the file is the interface

	var peerList []*Peer
	for i := 1; i < len(confElements); i++ { // parse rest of config file -- the rest are peers
		peer := parsePeer(confElements[i])
		peerList = append(peerList, peer)
	}

	interfaceName := strings.TrimSuffix(configurationFileName, ".conf") // remove .conf to get pure interface name
	conf := NewConfiguration(interfaceName, name, ifaceAddress, ifaceListenPort, ifacePrivateKey, ifaceDNS, peerList)

	return conf
}

// Return name, address, listenPort, privateKey, dns in that order
func parseInterfaceInfo(iface string) (string,string,int,string,string) {
	var name string
	var address string
	var listenPort int
	var privateKey string
	var dns string
	var err error
	
	lines := strings.Split(iface, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		namePrefix := "# Name = "
		if strings.HasPrefix(line, namePrefix) {
			name = strings.TrimPrefix(line, namePrefix)
		}

		addressPrefix := "Address = "
		if strings.HasPrefix(line, addressPrefix) {
			address = strings.TrimPrefix(line, addressPrefix)
		}

		listenPortPrefix := "ListenPort = "
		if strings.HasPrefix(line, listenPortPrefix) {
			listenPortStr := strings.TrimPrefix(line, listenPortPrefix)
			listenPort, err = strconv.Atoi(listenPortStr)
			if err != nil {
				slog.Error("Failed to parse listening port")
				panic(err)
			}
		}

		privateKeyPrefix := "PrivateKey = "
		if strings.HasPrefix(line, privateKeyPrefix) {
			privateKey = strings.TrimPrefix(line, privateKeyPrefix)
		}

		dnsPrefix := "DNS = "
		if strings.HasPrefix(line, dnsPrefix) {
			dns = strings.TrimPrefix(line, dnsPrefix)
		}

	}
	return name, address, listenPort, privateKey, dns
}

// Return name, publicKey, allowedIPs in that order
func parsePeer(peer string) *Peer {
	var name string
	var publicKey string
	var allowedIPs []string

	lines := strings.Split(peer, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		namePrefix := "# Name = "
		if strings.HasPrefix(line, namePrefix) {
			name = strings.TrimPrefix(line, namePrefix)
		}

		publicKeyPrefix := "PublicKey = "
		if strings.HasPrefix(line, publicKeyPrefix) {
			publicKey = strings.TrimPrefix(line, publicKeyPrefix)
		}
	
		allowedIPsPrefix := "AllowedIPs = "
		if strings.HasPrefix(line, allowedIPsPrefix) {
			allowedIPsLong := strings.TrimPrefix(line, allowedIPsPrefix)
			allowedIPsList := strings.Split(allowedIPsLong, ",")

			for _ ,addr := range allowedIPsList {
				ip := strings.Split(addr,"/")
				allowedIPs = append(allowedIPs, ip[0]) //we don't need after the slash i.e 10.5.5.1/32 -> 10.5.5.1
			}

		}
	}

	newPeer := NewPeer(name,publicKey,allowedIPs)

	return newPeer
}

func ParsePeerInfo(peerInfo string) *PeerInfo {
	var endPoint string
	var latestHandshake string
	transfer := make(map[string]string)

	lines := strings.Split(peerInfo, "\n")
	publicKey := strings.TrimSpace(lines[0]) // after string splits, pub key will be first line

	for _, line := range lines {
		line = strings.TrimSpace(line)
		endPointPrefix := "endpoint: "
		if strings.HasPrefix(line, endPointPrefix) {
			endPoint = strings.TrimPrefix(line, endPointPrefix)
		}

		latestHandshakePrefix := "latest handshake: "
		if strings.HasPrefix(line, latestHandshakePrefix) {
			latestHandshake = strings.TrimPrefix(line, latestHandshakePrefix)
		}

		transferPrefix := "transfer: "
		if strings.HasPrefix(line, transferPrefix) {
			transferTrim := strings.TrimPrefix(line, transferPrefix)
			transferS := strings.Split(transferTrim, ",")

			transfer["Received"] = strings.TrimSpace(transferS[0])
			transfer["Sent"] = strings.TrimSpace(transferS[1])
		} 
	}

	if transfer["Sent"] == "" {
		transfer["Sent"] = "0 Mib sent"
	}
	if transfer["Received"] == "" {
		transfer["Received"] = "0 Mib received"
	}

	newPeerInfo := NewPeerInfo(publicKey, endPoint,latestHandshake,transfer)
	return newPeerInfo
}

// call wg show <interfaceName> --> update peerInfoData
func LoadPeerInfo(interfaceName string, confs *WireGuardConfigurations) {
	wgInfo := string(WgSpecific(interfaceName))
	wgInfoElements := strings.Split(wgInfo, "peer:") //split on "peer:"


	for i := 1; i < len(wgInfoElements); i++ { // parse rest of info -- the rest are peers
		peerInfo := ParsePeerInfo(wgInfoElements[i])


		confs.ConfMap[interfaceName].PeerMap[peerInfo.PublicKey].Info = peerInfo
		confs.ConfMap[interfaceName].PeerMap[peerInfo.PublicKey].setStatus()
	}
}