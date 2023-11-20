package api

import (
	//"encoding/json"
	//"fmt"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	//"github.com/jcocozza/wg-monitor/utils"
	"github.com/jcocozza/wg-monitor/utils"
	s "github.com/jcocozza/wg-monitor/wireguard/structs"
	"github.com/jcocozza/wg-monitor/wireguard/commands"
)
type WgConfig map[string]*s.Configuration

func UpdateConfiguration(confs WgConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		configurationName := c.Param("configurationName")
		slog.Info("[API] Update Configuration: "+configurationName)
		// the configuration name is the wrong thing???
		if conf, ok := confs[configurationName]; ok && conf != nil {
			conf.Refresh()
			c.JSON(http.StatusOK, conf.Peers) //return the desired interface data from <confName>
		} else {
			// Handle the case where the configuration does not exist or is nil
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		}
	}
}

// this function will replate UpdateConfigurations()
func UpdateNetworks(confs WgConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusMap := make(map[string]bool)
		for _,conf := range confs {
			conf.NetworkInfo.CheckStatus()
			statusMap[conf.ConfName] = conf.NetworkInfo.Status
		}
		c.JSON(http.StatusOK, statusMap)
	}
}

// CombinedData represents the combined data structure
type NewPeerData struct {
	TextData    string `json:"textData"`
	QRCodeData  string `json:"qrCodeData"`
}

func AddPeer(wireguardPath string, confs WgConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		confName := c.Param("confName")
		confFilePath := wireguardPath+confName+".conf"
	
		nickName := c.PostForm("name")
		
		allowedIPsString := c.PostForm("allowedIPs")
		allowedIPs := strings.Split(allowedIPsString, ",")
		
		dns := c.PostForm("dns")
		
		vpnEndpoint := c.PostForm("vpnEndpoint")

		addressesToUseString := c.PostForm("addressesToUse")
		addressesToUse := strings.Split(addressesToUseString, ",")

		persistentKeepAliveString := c.PostForm("persistentKeepAlive")
		persistentKeepAlive, _ := strconv.Atoi(persistentKeepAliveString) // html form ensures that we get an integer
		
		peerFile, peerServerConf := confs[confName].GenerateNewPeer(confFilePath, nickName, allowedIPs, dns, vpnEndpoint, addressesToUse, persistentKeepAlive)

		/*
		data := map[string]interface{}{
			"name": name,
			"allowedIPs":allowedIPs,
			"dns":dns,
			"vpnEndpoint":vpnEndpoint,
			"addressToUse":addressesToUse,
			"persistentKeepAlive":persistentKeepAlive,
		}
		

		c.JSON(http.StatusOK, data)
		*/
		/*
		fmt.Println("Success:")
		fmt.Println("name:",confName)
		fmt.Println("name:",name)
		fmt.Println("name:",allowedIPs)
		fmt.Println("name:",dns)
		fmt.Println("name:",vpnEndpoint)
		fmt.Println("name:",addressesToUse)
		fmt.Println("name:",persistentKeepAlive)
		*/
		qrPeerData := utils.QRCodeData(peerFile, 300)
		str := base64.StdEncoding.EncodeToString(qrPeerData)
		data := NewPeerData{
			TextData: string(peerFile),
			QRCodeData: str,
		}

		c.JSON(http.StatusOK, data)
		//c.Header("Content-Type", "text/plain")
		//c.String(http.StatusOK, string(peerFile))
		commands.WgReloadServer(confName)
		utils.AppendTo(confFilePath, peerServerConf)
		confs[confName] = s.LoadConfiguration(confFilePath, confName)
	}
}

/* API call that also takes in an object
type Obj struct {}

func exampleAPICALL(obj *Obj) func(c *gin.Context) {
    return func(c *gin.Context) {
        // Your handling logic here using the serverInfo object and the Gin context
        // You can access serverInfo fields and perform necessary operations using the context
    }
}
*/