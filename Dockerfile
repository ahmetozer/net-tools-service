FROM ubuntu
RUN export DEBIAN_FRONTEND=noninteractive; apt update; apt install iputils-ping traceroute dnsutils mtr --no-install-recommends -y

FROM golang
WORKDIR /src
COPY . /src

FROM busybox:glibc
RUN rm -rf /bin/*

# COPY Binaries

COPY --from=0 /usr/bin/traceroute /bin/traceroute
COPY --from=0 /usr/bin/ping       /bin/ping
COPY --from=0 /usr/bin/mtr        /bin/mtr

#ENTRYPOINT ["/bin/traceroute"]
