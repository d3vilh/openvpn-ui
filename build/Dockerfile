# armv5l raspberry pi 1/zero 32bit: FROM arm32v5/debian:stable-slim, armv7l raspberry pi 2/3/4 32bit: FROM arm32v7/debian:stable-slim
# AMD64 64bit: FROM debian:stable-slim, arm64v8 64bit: FROM arm64v8/debian:stable-slim

FROM alpine:3.18.3
LABEL maintainer="Mr.Philipp <d3vilh@github.com>"
LABEL version="0.9"
WORKDIR /opt
EXPOSE 8080/tcp

RUN apk update \
 && apk add --update-cache \
    easy-rsa \
    curl \
    jq \
 && rm -rf /var/cache/apk/* && chmod 755 /usr/share/easy-rsa/*

COPY assets/*.sh /opt/scripts/
ADD openvpn-ui.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-gui/conf/app.conf

CMD /opt/start.sh