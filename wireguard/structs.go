package wireguard

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jcocozza/wg-monitor/utils"
)

type WireGuardConfigurations struct {
	ConfMap  map[string]*Configuration	`json:"confMap"`
	ConfNames []string					`json:"confNames"`
	////ConfMapPublicKey  map[string]*Configuration	`json:"confMapPublicKey"`
}

func NewWireGuardConfigurations(confList []*Configuration) *WireGuardConfigurations{
	var confNames []string
	confMap := make(map[string]*Configuration)
	for _, conf := range confList {
		confNames = append(confNames, conf.InterfaceName)
		confMap[conf.InterfaceName] = conf
		////ConfMapPublicKey[conf.PublicKey] = conf

	}

	return &WireGuardConfigurations{
		ConfMap: confMap,
		ConfNames: confNames,
	}
}

type Configuration struct {
	InterfaceName string			`json:"interfaceName"`
	Name          string        	`json:"name"`
	Address       string        	`json:"address"`
	ListenPort    int           	`json:"listenPort"`
	PrivateKey    string        	`json:"privateKey"`
	PublicKey     string        	`json:"publicKey"`
	DNS           string        	`json:"dns"`
	PostUp		  string			`json:"postUp"`
	PostDown	  string			`json:"postDown"`
	Peers         []*Peer 	    	`json:"peers"`
	PeerMap       map[string]*Peer	`json:"peerMap"`
}
func NewConfiguration(interfaceName string, name string, address string, listenPort int, privateKey string, dns string, postUp string, postDown string, peers []*Peer) *Configuration {
	if name == "" {
		name = interfaceName
	}
	publicKey := string(wgPubKey(privateKey))

	peerMap := make(map[string]*Peer)
	for _,peer := range peers {
		peerMap[peer.PublicKey] = peer
	}

	return &Configuration{
		InterfaceName: interfaceName,
		Name: name,
		Address: address,
		ListenPort: listenPort,	
		PrivateKey: privateKey,
		PublicKey: publicKey,
		DNS: dns,
		PostUp: postUp,
		PostDown: postDown,
		Peers: peers,
		PeerMap: peerMap,
	}
}

// Return how the server is represented in the server's .conf file.
// i.e. the stuff corresponding to [Interface] in the .conf of the server
// this is the file on the server machine
func (conf *Configuration) ServerConfFileOut() []byte {
	out := "[Interface]\n"
	out += fmt.Sprintf("# Name = %s\n", conf.Name)
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

type Peer struct {
	Name 		string			`json:"Name"`
	PublicKey 	string			`json:"PublicKey"`
	PrivateKey  string			`json:"PrivateKey"`
	AllowedIPs 	[]string		`json:"AllowedIPs"`
	Info 		*PeerInfo		`json:"Info"`
}

func NewPeer(name string, publicKey string, privateKey string, allowedIPs []string) *Peer {
	return &Peer{
		Name: name,
		PublicKey: publicKey,
		PrivateKey: privateKey,
		AllowedIPs: allowedIPs,
	}
	//// later possibly add a "getInfo" to finish class instatiation
}

type PeerInfo struct {
	PublicKey		string				`json:"PublicKey"`
	EndPoint 		string				`json:"EndPoint"`
	LatestHandshake string				`json:"LatestHandshake"`
	Transfer 		map[string]string	`json:"Transfer"` // transfer: sent/received 
	Status			bool				`json:"Status"`
}

func NewPeerInfo(publicKey string, endpoint string, latestHandshake string, transfer map[string]string) *PeerInfo {
	return &PeerInfo{
		PublicKey: publicKey,
		EndPoint: endpoint,
		LatestHandshake: latestHandshake,
		Transfer: transfer,
		Status: false,
	}
}

// check each of the peers allowed ips
func (p *Peer) setStatus() {
	totalNum := len(p.AllowedIPs)
	for i, addr := range p.AllowedIPs {
		slog.Debug(fmt.Sprintf("[%d/%d]: Checking %s...", i, totalNum, addr))
		p.Info.Status = utils.IsReachable(addr)
	}
}

// return how the peer is represented in the server's .conf file
// i.e the stuff corresponding to [Peer] in the .conf of the server
func (p *Peer) ServerConfFileOut() []byte {
	allowedIPsString := strings.Join(p.AllowedIPs,",")
	
	out := ("[Peer]\n")
	out += fmt.Sprintf("# Name = %s\n",p.Name)
	out += fmt.Sprintf("PublicKey = %s\n",p.PublicKey)
	out += fmt.Sprintf("AllowedIPs = %s\n", allowedIPsString)

	return []byte(out)
}

// generate the configuration file for a peer
// this is the file that will be on a CLIENT machine
func (p *Peer) confFileOut(dns string, vpnEndPoint string, conf *Configuration, addressesToUse []string, persistenKeepAlive int) []byte {
	allowedIPsString := strings.Join(p.AllowedIPs,",")
	addressesToUseString := strings.Join(addressesToUse,",")
	
	out := "[Interface]\n"
	if p.Name != "" {
		out += fmt.Sprintf("# Name = %s\n", p.Name)
	}
	out += fmt.Sprintf("Address = %s\n", allowedIPsString)
	out += fmt.Sprintf("PrivateKey = %s\n", p.PrivateKey)
	if dns != "" {
		out += fmt.Sprintf("DNS = %s\n", dns)
	}
	out += "[Peer]\n"
	out += fmt.Sprintf("# Name = %s\n",conf.Name)
	out += fmt.Sprintf("Endpoint = %s\n", vpnEndPoint)
	out += fmt.Sprintf("PublicKey = %s\n", conf.PublicKey)
	out += fmt.Sprintf("AllowedIPs = %s\n",addressesToUseString)

	if persistenKeepAlive == -999 {
		out += fmt.Sprintf("PersistentKeepalive = %d", persistenKeepAlive)
	}
	return []byte(out)
}
