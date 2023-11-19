package wireguard

import (
	"log/slog"
	"strings"

	c "github.com/jcocozza/wg-monitor/wireguard/commands"
	s "github.com/jcocozza/wg-monitor/wireguard/structs"
)

/*
0) user input - where is the wireguard folder located
1) Get relevant network interfaces
2) Get all wireguard .conf files
3) Create a configuration object for each
4) load peers
*/
func LoadWireGuard(wireguardPath string) map[string]*s.Configuration {

	networkInterfaces := s.LoadAllNetworkInterfaceInfo()

	files := strings.TrimSpace(string(c.GetConfNames(wireguardPath)))
	fileList := strings.Split(files, "\n")

	var confPath string
	configurationMap := make(map[string]*s.Configuration)

	for _, fileName := range fileList {
		slog.Debug("parsing conf:"+fileName)

		confPath = wireguardPath+"/"+fileName
		confName := strings.TrimSuffix(fileName,".conf")
		conf := s.LoadConfiguration(confPath, confName)

		if networkInterfaces[conf.PublicKey] != nil {
			conf.AttachNetwork(networkInterfaces[conf.PublicKey])
			conf.Refresh()
		} else {
			conf.AttachNetwork(s.EmptyNetworkInterface())
			conf.Refresh()
		}
		configurationMap[conf.ConfName] = conf
	}

	return configurationMap
}