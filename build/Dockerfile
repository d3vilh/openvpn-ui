# armv7l raspberry pi 2/3/4 32bit
#FROM arm32v7/debian:stable-slim

# armv5l raspberry pi 1/zero 32bit
#FROM arm32v5/debian:stable-slim 

# amd64 64bit
#FROM debian:stable-slim

# arm64v8 64bit
#FROM arm64v8/debian:stable-slim

FROM DEFINE-YOUR-ARCH
LABEL maintainer="Mr.Philipp <d3vilh@github.com>"
WORKDIR /opt
EXPOSE 8080

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y easy-rsa curl && apt-get clean && rm -rf /var/lib/apt/lists/* && chmod 755 /usr/share/easy-rsa/*
ADD assets/start.sh /opt/start.sh
ADD assets/generate_ca_and_server_certs.sh /opt/scripts/generate_ca_and_server_certs.sh
ADD assets/vars.template /opt/scripts/
ADD assets/genclient.sh /opt/scripts/genclient.sh
ADD assets/rmclient.sh /opt/scripts/rmclient.sh
ADD assets/restart.sh /opt/scripts/restart.sh
ADD assets/rmcert.sh /opt/scripts/rmcert.sh
ADD openvpn-ui.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-gui/conf/app.conf

CMD /opt/start.sh