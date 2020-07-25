FROM golang
WORKDIR /go/src/github.com/ahmetozer/net-tools-service
COPY . .

# Get dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# build
RUN CGO_ENABLED=0 go build -o /bin/net-tools-service

# Switch to server ENV
FROM ubuntu
# Get all required packages
RUN export DEBIAN_FRONTEND=noninteractive; apt update; apt install iputils-ping traceroute dnsutils whois ca-certificates curl mtr-tiny --no-install-recommends -y ; apt clean;apt autoclean

# COPY Binary and service file
COPY --from=0 /bin/net-tools-service        /bin/net-tools-service
COPY services /etc/services

# To allow net-tools-service to bind port root ports like a 443
RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/net-tools-service

# switch user to www-data
USER www-data
ENTRYPOINT ["/bin/net-tools-service"]