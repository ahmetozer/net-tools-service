# Net Tools Service

This software is backend service for [Looking Glass](https://github.com/ahmetozer/looking-glass) software.

Our Net Tools Service is a small Golang based tool to server ping, traceroute, nslookup, mtr, curl and whois tools.

It is designed to run basic net tools securely in containers.

## Features

- Automatic Self Certificate Generation
- Limit incoming requests (Rate Limit)
- Execute some network tools (nslookup, ping, tracert, whois, curl, mtr)
- Live command output
- Gracefully Shutdown
- Load Settings recursively from remote server with servers.json


## Installation

You can easily to deploy Looking Glass Service to your server with Docker.

```sh
docker run -it -p 443:443 ahmetozer.org/looking-glass-service --configURL https://lg.ahmetozer.org/server.json --svLoc Amsterdam.Amsterdam1
```

## Configuration

There is a few options to configure your software

- `--listenAddr` You can manage listen port and listen ip address with this argument.  
Ex. `--listenAddr 0.0.0.0:443`

- `--config-url` System loads server configs from frontend server which is example.com/server.json. This is prevent any front-end and back-end conflict.  
Ex. ` --config-url https://lg.ahmetozer.org/server.json `

- `--svloc` The system identifies itself with this argument  
Ex. `--svLoc Amsterdam.Amsterdam1`

- You can use also own SSL certificate to serve HTTPS server with verified certificate. To use own certificate mount your certificates to /cert/ folder in container.  
 Ex.  
  `docker run -it -p 443:443 --mount type=bind,source="/etc/letsencrypt/live/example.com/fullchain.pem",target=/cert/cert.pem,readonly  --mount type=bind,source="/etc/letsencrypt/live/example.com/privkey.pem",target=/cert/key.pem,readonly ahmetozer.org/looking-glass-service --config-url https://lg.ahmetozer.org/server.json --svloc Netherland.Amsterdam2`

### Serving to Global

If you have a extra IPv4 or IPv6, you can bind ips to container and server directly.  
If you don't have a extra IP, You can use nginx to serve multiple domains in one domain.  
Another way is bind container`s port to your server port. You can use any un used port expose service. Front end also support other ports.

### server.json

For server.json configuration please visit [https://github.com/ahmetozer/looking-glass#serverjson](https://github.com/ahmetozer/looking-glass#serverjson).

## Some Information's About this Software

- This program is only available for linux.

- Program is only executable in www-data user. (Prevent any security issue)

- If you use outside of container you have to install "ping", "traceroute", "whois", "nslookup", "mtr", "curl" into to your server.

- System is restrictive to better security. Please be check frontend website is added to referrers,
other wise frontend requests is blocked. Also be sure `Url` is right. System is automatically blocks other domains.
