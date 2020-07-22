FROM golang
WORKDIR /go/src/github.com/ahmetozer/looking-glass-service
COPY . .

# Get dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# build
RUN CGO_ENABLED=0 go build -o /bin/looking-glass-service

# Switch to server ENV
FROM ubuntu
# Get all required packages
RUN export DEBIAN_FRONTEND=noninteractive; apt update; apt install iputils-ping traceroute dnsutils whois ca-certificates curl --no-install-recommends -y ; apt clean;apt autoclean

# COPY Binary and service file
COPY --from=0 /bin/looking-glass-service        /bin/looking-glass-service
COPY services /etc/services

# To allow looking-glass-service to bind port root ports like a 443
RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/looking-glass-service

# switch user to www-data
USER www-data
ENTRYPOINT ["/bin/looking-glass-service"]