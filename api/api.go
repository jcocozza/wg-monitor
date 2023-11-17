package api

import (
	//"encoding/json"
	//"fmt"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	//"github.com/jcocozza/wg-monitor/utils"
	"github.com/jcocozza/wg-monitor/utils"
	"github.com/jcocozza/wg-monitor/wireguard"
)

func UpdateConfiguration(confs *wireguard.WireGuardConfigurations) func(c *gin.Context) {
	return func(c *gin.Context) {
		interfaceName := c.Param("interfaceName")

		wireguard.LoadPeerInfo(interfaceName,confs)

		/*
		jsonData, err := json.Marshal(confs.ConfMap[interfaceName])

		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}*/
		c.JSON(http.StatusOK, confs.ConfMap[interfaceName].Peers) //return the desired interface data from <interfaceName>
	}
}

// CombinedData represents the combined data structure
type NewPeerData struct {
	TextData    string `json:"textData"`
	QRCodeData  string `json:"qrCodeData"`
}

func AddPeer(wireguardPath string, confs *wireguard.WireGuardConfigurations) func(c *gin.Context) {
	return func(c *gin.Context) {
		interfaceName := c.Param("interfaceName")
		confFilePath := wireguardPath+interfaceName+".conf"
	
		name := c.PostForm("name")
		
		allowedIPsString := c.PostForm("allowedIPs")
		allowedIPs := strings.Split(allowedIPsString, ",")
		
		dns := c.PostForm("dns")
		
		vpnEndpoint := c.PostForm("vpnEndpoint")

		addressesToUseString := c.PostForm("addressesToUse")
		addressesToUse := strings.Split(addressesToUseString, ",")

		persistentKeepAliveString := c.PostForm("persistentKeepAlive")
		persistentKeepAlive, _ := strconv.Atoi(persistentKeepAliveString) // html form ensures that we get an integer
		
		

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
		fmt.Println("name:",interfaceName)
		fmt.Println("name:",name)
		fmt.Println("name:",allowedIPs)
		fmt.Println("name:",dns)
		fmt.Println("name:",vpnEndpoint)
		fmt.Println("name:",addressesToUse)
		fmt.Println("name:",persistentKeepAlive)
		*/
		peerFile := wireguard.GenerateNewPeer(confFilePath, name, allowedIPs, dns, vpnEndpoint, interfaceName, confs, addressesToUse, persistentKeepAlive)
		qrPeerData := utils.QRCodeData(peerFile, 300)
		str := base64.StdEncoding.EncodeToString(qrPeerData)
		data := NewPeerData{
			TextData: string(peerFile),
			QRCodeData: str,
		}

		c.JSON(http.StatusOK, data)
		//c.Header("Content-Type", "text/plain")
		//c.String(http.StatusOK, string(peerFile))
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