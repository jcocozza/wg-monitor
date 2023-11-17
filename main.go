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

func initWGMonitor(wireguardPath string) *wireguard.WireGuardConfigurations {
    wgConfs := wireguard.InitWireGuardConfigurations(wireguardPath)

    return wgConfs
}

// middleware for navlinks based on the configurations that are set up 
func SetNavLinks(confNames []string) gin.HandlerFunc {
    return func (c *gin.Context) {
        var navLinks []NavLink

        for _, confName := range confNames {
            navLinks = append(navLinks, NavLink{confName, fmt.Sprintf("/configurations/%s", confName)})
        }
        c.Set("navLinks", navLinks)
    }
}

func main() {
    wireguardPath := "/usr/local/etc/wireguard/" //os.Args[1]

    wgConfs := initWGMonitor(wireguardPath)

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

    router.GET("/configurations/:confName", func (c *gin.Context)  {
        confName := c.Param("confName")
        iface := wgConfs.ConfMap[confName]
        //iface := wireguardOld.ExtractInterface(interfaces, confName)

        //for _,peer := range iface.Peers{
        //    peer.CheckMetaStatus()
        //}

        c.HTML(http.StatusOK, "configuration.html", gin.H{
            "confName" : confName,
            "interface" : iface,
            "navLinks" : c.MustGet("navLinks").([]NavLink),
        })
    })


    router.GET("/configurations/:confName/newPeer", func(c *gin.Context) {
        c.HTML(http.StatusOK, "newPeerPopup.html", gin.H{})
    })

    // API ROUTES
    router.GET("/api/update/configurations/:confName", api.UpdateConfiguration(wgConfs))
    router.POST("/api/configurations/:confName/newPeer", api.AddPeer(wireguardPath, wgConfs))
    
    // Run the server
    //router.Run("143.229.244.67:8080")
    router.Run("10.5.5.1:8080")

}
