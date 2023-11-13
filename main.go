package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcocozza/wg-monitor/api"
	"github.com/jcocozza/wg-monitor/wireguard"
)

type NavLink struct {
    Text string
    URL  string
}

var templateVariable string // The variable that will be updated


// middleware for navlinks based on the configurations that are set up 
func SetNavLinks() gin.HandlerFunc {
    return func (c *gin.Context) {
        var navLinks []NavLink

        interfaceNames := wireguard.GetInterfaceNames()

        for _, interfaceName := range interfaceNames {
            navLinks = append(navLinks, NavLink{interfaceName, fmt.Sprintf("/configurations/%s", interfaceName)})
        }
        c.Set("navLinks", navLinks)
    }
}

func main() {
    // Initialize the Gin router
    router := gin.Default()
    router.Use(SetNavLinks())
    
    router.LoadHTMLGlob("web/templates/*")
    router.Static("/static", "web/static")

    interfaces := wireguard.GetWgInfo()

    // PAGES
    router.GET("/", func(c *gin.Context) {

        c.HTML(http.StatusOK, "home.html", gin.H{
            "interfaces" : interfaces,
            "navLinks" : c.MustGet("navLinks").([]NavLink),
        })
    })

    router.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "navLinks" : c.MustGet("navLinks").([]NavLink),
        })
    })

    router.GET("/configurations/:interfaceName", func (c *gin.Context)  {
        interfaceName := c.Param("interfaceName")

        iface := wireguard.ExtractInterface(interfaces, interfaceName)

        //for _,peer := range iface.Peers{
        //    peer.CheckMetaStatus()
        //}

        c.HTML(http.StatusOK, "configuration.html", gin.H{
            "interfaceName" : interfaceName,
            "interface" : iface,
            "navLinks" : c.MustGet("navLinks").([]NavLink),
        })
    })

    // API ROUTES
    router.GET("/api/getInterfaces", api.GetWireGuardServerInfo)
    router.GET("/api/configurations/:interfaceName", api.CheckPeerMetaStatus)
    
    // Run the server
    router.Run("10.5.5.1:8080")

}
