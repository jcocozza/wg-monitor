# wg-monitor
A monitoring dashboard for Wireguard - runs on Gin


## Assumptions

- wg-montior assumes that the name of the interface that your server is running on is the name of the .conf file.
    
    That is to say, `ifconfig <fileNamePrefix>` should return the interface corresponding to your server.
    This is probably only a problem of MacOS users. With MacOS you could name your config `wg0.conf`, but the interface the server runs on is `utunX`. So in this case, just ensure that you name it `utunX.conf`

- There is at most 1 PostUp. If you need to do more then 1 thing, simply put it all in a file and use the file path.
- There is at most 1 PostDown. If you need to do more then 1 thing, simply put it all in a file and use the file path.