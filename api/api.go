package api

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		
		
		fmt.Println("Success:")
		fmt.Println("name:",interfaceName)
		fmt.Println("name:",name)
		fmt.Println("name:",allowedIPs)
		fmt.Println("name:",dns)
		fmt.Println("name:",vpnEndpoint)
		fmt.Println("name:",addressesToUse)
		fmt.Println("name:",persistentKeepAlive)

		wireguard.GenerateNewPeer(confFilePath, name, allowedIPs, dns, vpnEndpoint, interfaceName, confs, addressesToUse, persistentKeepAlive)
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