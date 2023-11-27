# wg-monitor
A monitoring dashboard for Wireguard - runs on Gin

## Start Up

wg-monitor needs to know where your wireguard `.conf` files are stored. First, it will check the environment variable `$WIREGUARD_PATH`, then it check the first arguement passed into the script. Otherwise, it will use the default path: `/usr/local/etc/wireguard/`. 

## Assumptions
- There is at most 1 PostUp. If you need to do more then 1 thing, simply put it all in a file and use the file path.
- There is at most 1 PostDown. If you need to do more then 1 thing, simply put it all in a file and use the file path.