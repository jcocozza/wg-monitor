package main

import (
    "os"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcocozza/wg-monitor/api"
	"github.com/jcocozza/wg-monitor/wireguard"
	"github.com/jcocozza/wg-monitor/wireguard/structs"
)

type NavLink struct {
    Text string
    URL  string
    Status *structs.NetworkInterface
}


func initWGMonitor(wireguardPath string) api.WgConfig {
    wgConfs := wireguard.LoadWireGuard(wireguardPath)
    return wgConfs
}

// middleware for navlinks based on the configurations that are set up 
func SetNavLinks(configurations api.WgConfig) gin.HandlerFunc {
    return func (c *gin.Context) {
        var navLinks []NavLink

        for confName, conf := range configurations {
            navLinks = append(navLinks, NavLink{confName, fmt.Sprintf("/configurations/%s", confName), conf.NetworkInfo})
        }
        c.Set("navLinks", navLinks)
    }
}

func main() {

    var wireguardPath string
    var pathExists bool
    wireguardPath, pathExists = os.LookupEnv("WIREGUARD_PATH")

    if !pathExists {
        if len(os.Args) > 1 {
            wireguardPath = os.Args[1]
        } else {
            // default wireguard path
            wireguardPath = "/usr/local/etc/wireguard/"
        }
    }

    wgConfs := initWGMonitor(wireguardPath)

    // Initialize the Gin router
    router := gin.Default()
    router.Use(SetNavLinks(wgConfs))
    
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
        configuration := wgConfs[confName]
        
        c.HTML(http.StatusOK, "configuration.html", gin.H{
            "confName" : confName,
            "configuration" : configuration,
            "navLinks" : c.MustGet("navLinks").([]NavLink),
        })
    })


    router.GET("/configurations/:confName/newPeer", func(c *gin.Context) {
        c.HTML(http.StatusOK, "newPeerPopup.html", gin.H{})
    })

    // API ROUTES
    router.GET("/api/update/configurations/:configurationName", api.UpdateConfiguration(wgConfs))
    router.GET("/api/update/networks/all", api.UpdateNetworks(wgConfs))
    router.POST("/api/configurations/:confName/newPeer", api.AddPeer(wireguardPath, wgConfs))
    
    // Run the server
    //router.Run("143.229.244.67:8080")
    router.Run("10.5.5.1:8080")

}
