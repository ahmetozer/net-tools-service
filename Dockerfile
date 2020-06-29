FROM golang
WORKDIR /go/src/github.com/ahmetozer/looking-glass-controller
COPY . /go/src/github.com/ahmetozer/looking-glass-controller
RUN go get ./..
RUN CGO_ENABLED=0 go build -o /bin/looking-glass-controller
#RUN go build -o /bin/looking-glass-controller

FROM ubuntu
RUN export DEBIAN_FRONTEND=noninteractive; apt update; apt install iputils-ping traceroute dnsutils whois --no-install-recommends -y ; apt clean;apt autoclean

# COPY Binary and service file
COPY --from=0 /bin/looking-glass-controller        /bin/looking-glass-controller
COPY services /etc/services
# To allow looking-glass-controller to bind port 443
RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/looking-glass-controller

USER www-data
ENTRYPOINT ["/bin/looking-glass-controller"]
#CMD /bin/looking-glass-controller
