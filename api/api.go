package api

import (
	//"encoding/json"
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

/* API call that also takes in an object
type Obj struct {}

func exampleAPICALL(obj *Obj) func(c *gin.Context) {
    return func(c *gin.Context) {
        // Your handling logic here using the serverInfo object and the Gin context
        // You can access serverInfo fields and perform necessary operations using the context
    }
}
*/