package structs

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jcocozza/wg-monitor/utils"
)
type Peer struct {
	NickName 	string			`json:"name"`
	PublicKey 	string			`json:"publicKey"`
	PrivateKey  string			`json:"privateKey"`
	AllowedIPs 	[]string		`json:"allowedIPs"`
	Info 		*PeerInfo		`json:"info"`
	Parent      *Configuration	`json:"-"`
}

func NewPeer(nickName string, publicKey string, privateKey string, allowedIPs []string, parent *Configuration) *Peer {
	return &Peer{
		NickName: nickName,
		PublicKey: publicKey,
		PrivateKey: privateKey,
		AllowedIPs: allowedIPs,
		Info: EmptyPeerInfo(),
		Parent: parent,
	}
}

func (peer *Peer) UpdatePeerInfo(info *PeerInfo) {
	peer.Info = info
}

func (peer *Peer) SetPrivateKey(privateKey string) {
	peer.PrivateKey = privateKey
}

func (peer *Peer) SetParent(parent *Configuration) {
	peer.Parent = parent
	slog.Info("[Peer "+peer.PublicKey+"] Set Parent "+parent.ConfName)
}

// check each of the peer's allowed ips
func (peer *Peer) SetStatus() {
	totalNum := len(peer.AllowedIPs)
	if totalNum >= 1 {
		for i, addr := range peer.AllowedIPs {
			slog.Debug(fmt.Sprintf("[%d/%d]: Checking %s...", i, totalNum, addr))
			peer.Info.Status = utils.IsReachable(addr)
		}
	} else {
		peer.Info.Status = false
	}
}

// return how the peer is represented in the server's .conf file
// i.e the stuff corresponding to [Peer] in the .conf of the server
func (peer *Peer) ServerConfFileOut() []byte {
	allowedIPsString := strings.Join(peer.AllowedIPs,",")
	
	out := ("[Peer]\n")
	out += fmt.Sprintf("# Name = %s\n", peer.NickName)
	out += fmt.Sprintf("PublicKey = %s\n", peer.PublicKey)
	out += fmt.Sprintf("AllowedIPs = %s\n", allowedIPsString)

	return []byte(out)
}

// generate the configuration file for a peer
// this is the file that will be on a CLIENT machine
func (peer *Peer) ConfFileOut(dns string, vpnEndPoint string, addressesToUse []string, persistenKeepAlive int) []byte {
	allowedIPsString := strings.Join(peer.AllowedIPs,",")
	addressesToUseString := strings.Join(addressesToUse,",")
	
	out := "[Interface]\n"
	if peer.NickName != "" {
		out += fmt.Sprintf("# Name = %s\n", peer.NickName)
	}
	out += fmt.Sprintf("Address = %s\n", allowedIPsString)
	out += fmt.Sprintf("PrivateKey = %s\n", peer.PrivateKey)
	if dns != "" {
		out += fmt.Sprintf("DNS = %s\n", dns)
	}
	out += "[Peer]\n"
	out += fmt.Sprintf("# Name = %s\n", peer.Parent.NickName)
	out += fmt.Sprintf("Endpoint = %s\n", vpnEndPoint)
	out += fmt.Sprintf("PublicKey = %s\n", peer.Parent.PublicKey)
	out += fmt.Sprintf("AllowedIPs = %s\n",addressesToUseString)

	if persistenKeepAlive == -999 {
		out += fmt.Sprintf("PersistentKeepalive = %d", persistenKeepAlive)
	}
	return []byte(out)
}


type PeerInfo struct {
	PublicKey		string				`json:"publicKey"`
	EndPoint 		string				`json:"endPoint"`
	LatestHandshake string				`json:"latestHandshake"`
	Transfer 		map[string]string	`json:"transfer"` // transfer: sent/received 
	Status			bool				`json:"status"`
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

func EmptyPeerInfo() *PeerInfo {
	// Initialize peer info data
	transfer := make(map[string]string)
	transfer["Sent"] = "0 Mib sent"
	transfer["Received"] = "0 Mib sent"
	return &PeerInfo{
		PublicKey: "",
		EndPoint: "",
		LatestHandshake: "",
		Transfer: transfer,
		Status: false,
	}
}

