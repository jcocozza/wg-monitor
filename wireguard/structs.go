package wireguard

import (
	"fmt"
	"log/slog"

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
	Peers         []*Peer 	    	`json:"peers"`
	PeerMap       map[string]*Peer	`json:"peerMap"`
}
func NewConfiguration(interfaceName string, name string, address string, listenPort int, privateKey string, dns string, peers []*Peer) *Configuration {
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
		Peers: peers,
		PeerMap: peerMap,
	}
}

type Peer struct {
	Name 		string			`json:"Name"`
	PublicKey 	string			`json:"PublicKey"`
	AllowedIPs 	[]string		`json:"AllowedIPs"`
	Info 		*PeerInfo		`json:"Info"`
}

func NewPeer(name string, publicKey string, allowedIPs []string) *Peer {
	return &Peer{
		Name: name,
		PublicKey: publicKey,
		AllowedIPs: allowedIPs,
	}
	// later possibly add a "getInfo" to finish class instatiation
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