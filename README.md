# Net Tools Service

This software is backend service for [Looking Glass](https://github.com/ahmetozer/looking-glass) and [LatencyTable](https://github.com/ahmetozer/latencytable).

Our Net Tools Service is a small Golang based tool to serve ping, traceroute, nslookup, mtr, curl and whois tools.

It is designed to run basic net tools securely in containers.

## Features

- Automatic Self Certificate Generation
- Limit incoming requests (Rate Limit)
- Execute some network tools (nslookup, ping, tracert, whois, curl, mtr, svinfo)
- Live command output
- Gracefully Shutdown

## Installation

You can easily to deploy Net Tools Service to your server with Docker.

```sh
docker run -it -p 443:443 ghcr.io/ahmetozer/net-tools-service:latest
```

## Configuration

There is a few options to configure your software

### Environment variables

- **functions**  
To allow only few functions.
By default, if you don't give any function, all function is enabled.  
E.g.  `docker run -it -e functions="ping,whois"  ghcr.io/ahmetozer/net-tools-service:latest`

- **referrers**  
To allow only given referrers (Which website make request to this service).
By default, if you don't give any domain, all domains is accepted.  
E.g.  `docker run -it -e referrers="lg.ahmetozer.org"  ghcr.io/ahmetozer/net-tools-service:latest`  
`docker run -it -e referrers="lg.ahmetozer.org,noc.ahmetozer.org"  ghcr.io/ahmetozer/net-tools-service:latest`

- **hostname**  
To bind listen port to only one domain.
By default, all domains are accepted.  
E.g.  `docker run -it -e hostname="marmaris1-nts.ahmetozer.org"  ghcr.io/ahmetozer/net-tools-service:latest`  
If you give a empty variable, system uses containers hostname to bind domain.  
E.g.  `docker run -it -e hostname=""  ghcr.io/ahmetozer/net-tools-service:latest`  

- **ipver**  
To allow IP version on this system.
If it's empty, IPv4 and IPv6 are enabled. You can also enable IPv4 and IPv6 with DS or give IPv4 and IPv6 in env variable. If you define only IPv4 or IPv6, just defined IP version is enabled.  
E.g.  `docker run -it -e functions="ping,whois" -e ipver="IPv4" ghcr.io/ahmetozer/net-tools-service:latest`  
To disable IPv6 and IPv4 at same time (it means disable server), define ipver to "disabled"  
E.g.  `docker run -it -e functions="ping,whois" -e ipver="disabled" ghcr.io/ahmetozer/net-tools-service:latest`

- **rate**  
To change request limit in one second.  
Default is one request in one second.
If you require more request in one second (E.g. [latencytable](https://github.com/ahmetozer/latencytable) require)  
E.g.  `docker run -it -e rate="10" -e ipver="" ghcr.io/ahmetozer/net-tools-service:latest`

- **cache**  
tcp and icmp functions is cache module and cache is setted to 10second by default.  
E.g.  `docker run -it -e rate="10" -e cache="10s" ghcr.io/ahmetozer/net-tools-service:latest`

- **listenaddr**  
You can manage listen port and listen ip address with this argument.  
Mostly not used.  
E.g.  `docker run -it -e listenaddr="198.51.100.5:443" ghcr.io/ahmetozer/net-tools-service:latest`  
E.g.  `docker run -it -e listenaddr=":8443" ghcr.io/ahmetozer/net-tools-service:latest`

### Custom SSL Certificate

- You can use also own SSL certificate to serve HTTPS server with verified certificate. To use own certificate mount your certificates to /cert/ folder in container.  
 E.g.  

 ```bash
docker run -it -p 443:443 \
--mount type=bind,source="/etc/letsencrypt/live/example.com/fullchain.pem",target=/cert/cert.pem,readonly \
 --mount type=bind,source="/etc/letsencrypt/live/example.com/privkey.pem",target=/cert/key.pem,readonly \
 ghcr.io/ahmetozer/net-tools-service:latest
```

### Serving to Global

If you have a extra IPv4 or IPv6, you can bind ips to container and server directly.  
If you don't have a extra IP, You can use nginx to serve multiple domains in one domain.  
Another way is bind container`s port to your server port. You can use any un used port to expose net tools service.

## Some Information's About this Software

- This program is only available for linux.

- Program is only executable in www-data user. (Prevent any security issue)

- If you use outside of container you have to install "ping", "traceroute", "whois", "nslookup", "mtr", "curl" into to your server.

- You can get your server IP address with svinfo function.
