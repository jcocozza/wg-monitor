package wireguard

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-ping/ping"
)

/* The contents of `sudo wg show all dump`

interface server-private-key server-public-key listen-port

*/

// Check if a given address in up
func isReachable(addr string) bool {
	slog.Debug(fmt.Sprintf("pinging %s", addr))

	pinger, err := ping.NewPinger(addr)
	if err != nil {
		slog.Error("Error creating pinger:", err)
		return false
	}

	pinger.Count = 1 // Number of ICMP packets to send
	pinger.Timeout = 1 * time.Second // Timeout for each packet

	err = pinger.Run()
    if err != nil {
        slog.Error("Failed to ping:", err)
        return false
    }
	stats := pinger.Statistics()

	if stats.PacketsSent == stats.PacketsRecv {
		slog.Debug(fmt.Sprintf("%s is up...", addr))
	} else {
		slog.Debug(fmt.Sprintf("%s is down...", addr))
	}
	return stats.PacketsSent == stats.PacketsRecv
}

type Peer struct {
	PublicKey string
	EndPoint string
	AllowedIPs []string
	LatestHandshake string
	Transfer []string
	Status map[string]bool
	MetaStatus bool
}

func NewPeer(publicKey string, endPoint string, allowedIPs []string, latestHandshake string, transfer []string) *Peer {
	p := &Peer{
		PublicKey: publicKey,
		EndPoint: endPoint,
		AllowedIPs: allowedIPs,
		LatestHandshake: latestHandshake,
		Transfer: transfer,
		MetaStatus: false,
	}
	
	p.setStatus()
	return p
}

// check each of the peers allowed ips
func (p *Peer) setStatus() {
	totalNum := len(p.AllowedIPs)
	p.Status = make(map[string]bool)
	for i, addr := range p.AllowedIPs {
		slog.Debug(fmt.Sprintf("[%d/%d]: Checking %s...", i, totalNum, addr))
		p.Status[addr] = isReachable(addr)
	}
}

func (p *Peer) updateMetaStatus(status bool) {
    p.MetaStatus = status
}

// update status of peer, check if any IPs are up 
// update meta status,
// return true if at least 1 ip address is up
func (p *Peer) CheckMetaStatus() bool {
	//p.setStatus()

	for _, v := range p.Status {
		if v {
			p.updateMetaStatus(true)
			return true
		}
	}
	p.updateMetaStatus(false)
	return false
} 

type Interface struct {
	Name string
	IPAddress string
	PublicKey string
	ListeningPort string
	Peers []Peer
	PeerMap map[string]*Peer
}

// Initialize interface, starts with no peers
func NewInterface(name string, publickey string, listeningPort string) *Interface {
	peers := []Peer{}
	peerMap := make(map[string]*Peer)

	return &Interface{
		Name: name,
		IPAddress: GetInterfaceIP(name),
		PublicKey: publickey,
		ListeningPort: listeningPort,
		Peers: peers,
		PeerMap: peerMap,
	}
}

// Add a peer to the interface
func (iface *Interface) AddPeer(peer Peer) {
	iface.Peers = append(iface.Peers, peer)
	iface.BuildPeerMap()
}

// Construct a map of peer public keys to their corresponding peer object
func (iface *Interface) BuildPeerMap() {
	for _,peer := range iface.Peers {
		iface.PeerMap[peer.PublicKey] = &peer
	}
}

func ExtractInterface(list []Interface, criteria string) (Interface) {
	for _, element := range list {
		if element.Name == criteria {
			return element
		}
	}
	return Interface{}
}