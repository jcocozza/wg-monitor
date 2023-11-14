package api

import (
	//"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jcocozza/wg-monitor/wireguard"
)

func GetWireGuardServerInfo(c *gin.Context) {
	interfaces := wireguard.GetWgInfo()
	c.JSON(http.StatusOK, interfaces)
}


func CheckPeerMetaStatus(c *gin.Context) {
	interfaceName := c.Param("interfaceName")
	iface := wireguard.GetInterfaceByName(interfaceName)

	for i, _ := range iface.Peers{
		iface.Peers[i].CheckMetaStatus()
	}
	
	c.JSON(http.StatusOK, iface)
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