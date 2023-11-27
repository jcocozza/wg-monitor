# wg-monitor
A monitoring dashboard for Wireguard.

wg-monitor is database-less and runs on a single binary. 

Draws inspiration from https://github.com/donaldzou/WGDashboard and https://github.com/perara/wg-manager

The basic idea is a simple monitoring dashboard that runs on the same machine that the Wireguard VPN server is running on.

## Functionality
### Current capabilities:

- monitoring wireguard configurations
- adding peers to a configuration
- Turning configurations on/off (although...I'm not sure this is working 100% right now...)  
 
### Future goals/to-do:
- create/delete entire configurations from gui
- remove peers from a configuration
    - Allows for file download and QR code scanning
- make the front-end not so terrible 
    - pretty up html
    - code the javascript properly (maybe use something like htmx to remove it entirely?)   

## Start Up

wg-monitor needs to know where your wireguard `.conf` files are stored. First, it will check the environment variable `$WIREGUARD_PATH`,if that is empty, then it checks the first argument passed into the script. Otherwise, it will use the default path: `/usr/local/etc/wireguard/`. 

wg-monitor needs sudo access so it can run the `wg-quick` and `wg show` commands.
To run it simply do: `sudo ./<path/to/wg-monitor-binary>`. This will start up the dashboard on the machine and make it available on port `8080`. 

## Assumptions
- There is at most 1 PostUp. If you need to do more than 1 thing, simply put it all in a file and use the file path.
- There is at most 1 PostDown. If you need to do more than 1 thing, simply put it all in a file and use the file path.