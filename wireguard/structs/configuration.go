package structs

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jcocozza/wg-monitor/utils"
	c "github.com/jcocozza/wg-monitor/wireguard/commands"
)

type Configuration struct {
	ConfName      string			`json:"confName"`
	NickName      string        	`json:"nickName"`
	Address       string        	`json:"address"`
	ListenPort    int           	`json:"listenPort"`
	PrivateKey    string        	`json:"privateKey"`
	PublicKey     string        	`json:"publicKey"`
	DNS           string        	`json:"dns"`
	PostUp		  string			`json:"postUp"`
	PostDown	  string			`json:"postDown"`
	Peers         []*Peer 	    	`json:"peers"`
	PeerMap       map[string]*Peer	`json:"peerMap"`
	NetworkInfo   *NetworkInterface	`json:"network"`
}


func NewConfiguration(confName string, nickName string, address string, listenPort int, privateKey string, dns string, postUp string, postDown string) *Configuration {
	if nickName == "" {
		nickName = confName
	}
	publicKey := strings.TrimSpace(string(c.WgPubKey(privateKey)))

	peerMap := make(map[string]*Peer)

	return &Configuration{
		ConfName: confName,
		NickName: nickName,
		Address: address,
		ListenPort: listenPort,	
		PrivateKey: privateKey,
		PublicKey: publicKey,
		DNS: dns,
		PostUp: postUp,
		PostDown: postDown,
		PeerMap: peerMap,
		NetworkInfo: EmptyNetworkInterface(),
	}
}

// load in a wireguard configuration file e.g. `wg0.conf`
// this function is stupidly long
// Read configuration file, create a configuration object with peers
func LoadConfiguration(configurationPath string, confName string) *Configuration {
	file, err := utils.ReadFile(configurationPath)

	if err != nil {
		slog.Error("Failed to read configuration file")
	}

	conf := string(file)

	confElements := strings.Split(conf, "[Peer]") //split on "[Peer]"

	configurationInterface := confElements[0] // the first thing in the file is the interface
	configurationInterfaceLines := strings.Split(configurationInterface, "\n")

	var confNickName string
	var confAddress string
	var confListenPort int
	var confPrivateKey string
	var confDns string
	var confPostUp string
	var confPostDown string
	for _, ln := range configurationInterfaceLines {
		ln = strings.TrimSpace(ln)


		confNickNamePrefix := "# Name = "
		if strings.HasPrefix(ln, confNickNamePrefix) {
			confNickName = strings.TrimSpace(strings.TrimPrefix(ln, confNickNamePrefix))
			slog.Info("[Configuration " + confName + "]: Peer Nickname " + confNickName)
		}

		confAddressPrefix := "Address = "
		if strings.HasPrefix(ln, confAddressPrefix) {
			confAddress = strings.TrimSpace(strings.TrimPrefix(ln, confAddressPrefix))
			slog.Info("[Configuration " + confName + "]: Configuration Address " + confAddress)
		}

		confListenPortPrefix := "ListenPort = "
		if strings.HasPrefix(ln, confListenPortPrefix) {
			listenPortStr := strings.TrimPrefix(ln, confListenPortPrefix)
			confListenPort, err = strconv.Atoi(listenPortStr)
			if err != nil {
				slog.Error("Failed to parse listening port")
				panic(err)
			}
			slog.Info("[Configuration " + confName + "]: Configuration ListenPort " + listenPortStr)
		}

		confPrivateKeyPrefix := "PrivateKey = "
		if strings.HasPrefix(ln, confPrivateKeyPrefix) {
			confPrivateKey = strings.TrimSpace(strings.TrimPrefix(ln, confPrivateKeyPrefix))
			slog.Info("[Configuration " + confName + "]: Configuration Private Key " + confPrivateKey)
		}

		confDnsPrefix := "DNS = "
		if strings.HasPrefix(ln, confDnsPrefix) {
			confDns = strings.TrimSpace(strings.TrimPrefix(ln, confDnsPrefix))
			slog.Info("[Configuration " + confName + "]: Configuration DNS " + confDns)
		}

		confpostUpPrefix := "PostUp = "
		if strings.HasPrefix(ln, confpostUpPrefix) {
			confPostUp = strings.TrimSpace(strings.TrimPrefix(ln, confpostUpPrefix))
			slog.Info("[Configuration " + confName + "]: Configuration PostUp " + confPostUp)
		}

		confPostDownPrefix := "PostDown = "
		if strings.HasPrefix(ln, confPostDownPrefix) {
			confPostDown = strings.TrimSpace(strings.TrimPrefix(ln, confPostDownPrefix))
			slog.Info("[Configuration " + confName + "]: Configuration PostDown " + confPostDown)
		}

	}

	configuration := NewConfiguration(confName, confNickName, confAddress, confListenPort, confPrivateKey, confDns, confPostUp, confPostDown)


	var tempPeer *Peer
	var peerNickName string
	var peerPublicKey string
	var peerAllowedIPs []string
	for _, peerString := range confElements[1:] { // parse rest of config file -- the rest are peers	
		lines := strings.Split(peerString, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
	
			nickNamePrefix := "# Name = "
			if strings.HasPrefix(line, nickNamePrefix) {
				peerNickName = strings.TrimSpace(strings.TrimPrefix(line, nickNamePrefix))
				slog.Info("[Configuration " + confName + "][Peer]: Peer Nickname " + peerNickName)
			}
	
			publicKeyPrefix := "PublicKey = "
			if strings.HasPrefix(line, publicKeyPrefix) {
				peerPublicKey = strings.TrimSpace(strings.TrimPrefix(line, publicKeyPrefix))
				slog.Info("[Configuration " + confName + "][Peer]: Peer Public Key " + peerPublicKey)
			}
		
			allowedIPsPrefix := "AllowedIPs = "
			if strings.HasPrefix(line, allowedIPsPrefix) {
				allowedIPsLong := strings.TrimPrefix(line, allowedIPsPrefix)
				slog.Info("[Configuration " + confName + "][Peer]: Allowed IPs " + allowedIPsLong)
				allowedIPsList := strings.Split(allowedIPsLong, ",")
				for _ ,addr := range allowedIPsList {
					ip := strings.Split(addr,"/")
					peerAllowedIPs = append(peerAllowedIPs, ip[0]) //we don't need after the slash i.e. 10.5.5.1/32 -> 10.5.5.1
				}
	
			}
		}
		tempPeer = NewPeer(peerNickName, peerPublicKey,"",peerAllowedIPs, configuration)
		configuration.AddPeer(tempPeer)

		// clear vars
		peerNickName = ""
		peerPublicKey = ""
		peerAllowedIPs = []string{}
	}
	return configuration
}

// Attach network info to the configuration
func (conf *Configuration) AttachNetwork(network *NetworkInterface) {
	//slog.Info("Attaching Network: " + network.Name)
	conf.NetworkInfo = network
	slog.Info("[Configuration " + conf.ConfName + "] Attaching Network: " + conf.NetworkInfo.Name)
}

// Add a peer to the configuration
func (conf *Configuration) AddPeer(peer *Peer) {
	conf.PeerMap[peer.PublicKey] = peer
	conf.Peers = append(conf.Peers, peer)
	peer.SetParent(conf)
	slog.Info("[Configuration " + conf.ConfName + "] Added Peer: " + conf.PeerMap[peer.PublicKey].PublicKey)
}

// parse the output of wg show <interfaceName> to determine peer data
// another exceedingly long function
// Essentially: run wg show <interfaceName>, parse peer info data
func (conf *Configuration) Refresh() {
	if !conf.NetworkInfo.CheckStatus() { //if network is down
		return
	}

	wgInfo := strings.Split(string(c.WgSpecific(conf.NetworkInfo.Name)), "peer:") // split on "peer:"
	var lines []string
	var publicKey string
	var endPoint string
	var latestHandshake string
	var transfer map[string]string


	endPointPrefix := "endpoint: "
	latestHandshakePrefix := "latest handshake: "
	transferPrefix := "transfer: "
	for i := 1; i < len(wgInfo); i++ {
		peerInfo := wgInfo[i]
		lines = strings.Split(peerInfo, "\n")
		publicKey = strings.TrimSpace(lines[0]) // after string splits, pub key will be first line
		transfer = make(map[string]string)
		for _,line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, endPointPrefix) {
				endPoint = strings.TrimPrefix(line, endPointPrefix)
			}
			if strings.HasPrefix(line, latestHandshakePrefix) {
				latestHandshake = strings.TrimPrefix(line, latestHandshakePrefix)
				
			}
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

		slog.Info("[Configuration " + conf.ConfName + "][Peer "+publicKey+"][PeerInfo]: Endpoint " + endPoint)
		slog.Info("[Configuration " + conf.ConfName + "][Peer "+publicKey+"][PeerInfo]: Latest Handshake " + latestHandshake)
		slog.Info("[Configuration " + conf.ConfName + "][Peer "+publicKey+"][PeerInfo]: Received " + transfer["Received"])
		slog.Info("[Configuration " + conf.ConfName + "][Peer "+publicKey+"][PeerInfo]: Sent " + transfer["Sent"])

		conf.PeerMap[publicKey].UpdatePeerInfo(NewPeerInfo(publicKey, endPoint, latestHandshake, transfer))
		conf.PeerMap[publicKey].SetStatus() // after peer info has been updated, we can check to see if the peer is actually online

		endPoint = ""
		latestHandshake = ""
	}
}

func (conf *Configuration) GenerateNewPeer(configurationPath string, peerNickName string, allowedIPs []string, dns string, vpnEndpoint string, addressesToUse []string, persistentKeepAlive int) []byte{
	privateKey, publicKey := c.GenerateKeyPair()
	peer := NewPeer(peerNickName, publicKey, privateKey, allowedIPs, conf)
	out := peer.ConfFileOut(dns, vpnEndpoint, addressesToUse, persistentKeepAlive)

	/*
	sanatizedKey := utils.SanatizeKey(peer.PublicKey)
	peerFolderPath := "web/createdPeers/"+sanatizedKey
	utils.MkDir(peerFolderPath)
	utils.GenerateQRCode(peerFolderPath+"/qrcode.png", out)
	utils.WriteFile(peerFolderPath+"/peer.conf",out)
	utils.WriteFile(peerFolderPath+"/peerSever.conf",peer.ServerConfFileOut())
	*/
	//utils.AppendTo(configurationPath,peer.ServerConfFileOut()) add the peer to the server's conf file
	//ReloadServer(confName) reload the wireguard server in the background
	
	//"web/qrcodes/qrcode"+name+".png"
	return out
}

// Return how the server is represented in the server's .conf file.
// i.e. the stuff corresponding to [Interface] in the .conf of the server
// this is the file on the server machine
func (conf *Configuration) ServerConfFileOut() []byte {
	out := "[Interface]\n"
	out += fmt.Sprintf("# Name = %s\n", conf.NickName)
	out += fmt.Sprintf("Address = %s\n", conf.Address)
	out += fmt.Sprintf("ListenPort = %d\n", conf.ListenPort)
	out += fmt.Sprintf("PrivateKey = %s\n", conf.PrivateKey)
	out += fmt.Sprintf("DNS = %s\n", conf.DNS)
	out += fmt.Sprintf("PostUp = %s\n", conf.PostUp)
	out += fmt.Sprintf("PostDown = %s\n", conf.PostDown)
	// add in peers
	for _, peer := range conf.Peers {
		out += "\n"
		out += string(peer.ServerConfFileOut())
	}

	return []byte(out)
}

