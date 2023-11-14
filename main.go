package main

import (
	"fmt"
	//"log/slog"
	"net/http"

	//"os"

	"github.com/gin-gonic/gin"

	"github.com/jcocozza/wg-monitor/api"
	"github.com/jcocozza/wg-monitor/wireguard"
)

type NavLink struct {
    Text string
    URL  string
}

func initWGMonitor() *wireguard.WireGuardConfigurations {
    wireguardPath := "/usr/local/etc/wireguard/" //os.Args[1]
    wgConfs := wireguard.InitWireGuardConfigurations(wireguardPath)

    return wgConfs
}

// middleware for navlinks based on the configurations that are set up 
func SetNavLinks(confNames []string) gin.HandlerFunc {
    return func (c *gin.Context) {
        var navLinks []NavLink

        for _, interfaceName := range confNames {
            navLinks = append(navLinks, NavLink{interfaceName, fmt.Sprintf("/configurations/%s", interfaceName)})
        }
        c.Set("navLinks", navLinks)
    }
}

func main() {

    wgConfs := initWGMonitor()

    // Initialize the Gin router
    router := gin.Default()
    router.Use(SetNavLinks(wgConfs.ConfNames))
    
    router.LoadHTMLGlob("web/templates/*")
    router.Static("/static", "web/static")

    // PAGES
    router.GET("/", func(c *gin.Context) {

        c.HTML(http.StatusOK, "home.html", gin.H{
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
        iface := wgConfs.ConfMap[interfaceName]
        //iface := wireguardOld.ExtractInterface(interfaces, interfaceName)

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
    router.GET("/api/update/configurations/:interfaceName", api.UpdateConfiguration(wgConfs))
    
    // Run the server
    router.Run("10.5.5.1:8080")

}
