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