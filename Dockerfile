FROM golang
WORKDIR /go/src/github.com/ahmetozer/looking-glass-service
COPY . /go/src/github.com/ahmetozer/looking-glass-service
RUN go get ./..
RUN CGO_ENABLED=0 go build -o /bin/looking-glass-service
#RUN go build -o /bin/looking-glass-service

FROM ubuntu
RUN export DEBIAN_FRONTEND=noninteractive; apt update; apt install iputils-ping traceroute dnsutils whois --no-install-recommends -y ; apt clean;apt autoclean

# COPY Binary and service file
COPY --from=0 /bin/looking-glass-service        /bin/looking-glass-service
COPY services /etc/services
# To allow looking-glass-service to bind port 443
RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/looking-glass-service

USER www-data
ENTRYPOINT ["/bin/looking-glass-service"]