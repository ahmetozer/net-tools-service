# Looking Glass Service

This software is backend service for [Looking Glass](https://github.com/ahmetozer/looking-glass) software.

Our network looking glass service is a small Golang based tool to server ping, traceroute, nslookup and whois tools.

It is designed to run basic net tools securitly in containers to isolate from your server.

#### Features 

- Automatic Self Certifacate Generation
- Limit incoming requests (Rate Limit)
- Execute some network tools (nslookup, ping, tracert, whois)
- Live command output
- Gracefull Shutdown
- Load Settings recursively from servers.json

### Installation

You can easily to deploy Looking Glass Service to your server with Docker.

```sh
docker run -it -p 443:443 ahmetozer.org/looking-glass-service --configURL https://lg.ahmetozer.org/server.json --svLoc Amsterdam.Amsterdam1
```

#### Configuration

There is a few options to configure your software

- `--listenAddr` You can manage listen port and listen ip addres with this argumant.  
Ex. `--listenAddr 0.0.0.0:443`

- `--configURL` System loads server configs from frontend server which is example.com/server.json. This is prevent any front-end and back-end conflict.  
Ex. ` --configURL https://lg.ahmetozer.org/server.json `

- `--svLoc` The system identifies itself with this argument  
Ex. `--svLoc Amsterdam.Amsterdam1`

- You can use also own SSL certificate to serve HTTPS server with verified certificate. To use own certificate mount your certificates to /cert/ folder in container.  
 Ex.  
  `docker run -it -p 443:443 --mount type=bind,source="/etc/letsencrypt/live/example.com/fullchain.pem",target=/cert/cert.pem,readonly  --mount type=bind,source="/etc/letsencrypt/live/example.com/privkey.pem",target=/cert/key.pem,readonly ahmetozer.org/looking-glass-service --configURL https://lg.ahmetozer.org/server.json --svLoc Amsterdam.Amsterdam1 `

#### Serving to Global

If you have a extra IPv4 or IPv6, you can bind ips to container and server directly.  
If you dont have a extra IP, You can use nginx to serve multiple domains in one domain.  
Another way is bind container`s port to your server port. You can use any un used port expose service. Front end also support other ports.

##### Some Information's About this Software

- This program is only available for linux.

- Program is only executable in www-data user. (Prevent any security issue)

- If you use outside of container you have install "ping", "traceroute", "whois", "nslookup" commands to your server.

- You can`t separate server.json and WEBUI, If the given config url is different from WEBUI/server.json Program automatically block request to preventing Cross-Origin.