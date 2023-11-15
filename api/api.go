package api

import (
	//"encoding/json"
	"fmt"
	"net/http"

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

func NewPeer(c *gin.Context) {
	interfaceName := c.Param("interfaceName")
	
	name := c.PostForm("name")
    allowedIPs := c.PostForm("allowedIPs")
	dns := c.PostForm("dns")
    vpnEndpoint := c.PostForm("vpnEndpoint")
	addressesToUse := c.PostForm("addressesToUse")
	persistentKeepAlive := c.PostForm("persistentKeepAlive")

	fmt.Println("Success:")
	fmt.Println("name:",interfaceName)
	fmt.Println("name:",name)
	fmt.Println("name:",allowedIPs)
	fmt.Println("name:",dns)
	fmt.Println("name:",vpnEndpoint)
	fmt.Println("name:",addressesToUse)
	fmt.Println("name:",persistentKeepAlive)
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