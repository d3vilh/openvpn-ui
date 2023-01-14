# OpenVPN-web-ui

## Summary
OpenVPN server web administration interface.

Goal: create quick to deploy and easy to use solution that makes work with small OpenVPN environments a breeze.

If you have docker and docker-compose installed, you can jump directly to [installation](#Prod).

![Status page](docs/images/preview_status.png?raw=true)

Please note this project is in alpha stage. It still needs some work to make it secure and feature complete.

## Motivation



## Features

* status page that shows server statistics and list of connected clients
* easy creation of client certificates
* ability to download client certificates as a zip package with client configuration inside
* log preview
* modification of OpenVPN configuration file through web interface

## Screenshots

[Screenshots](docs/screenshots.md)

## Usage

After startup web service is visible on port 8080. To login use the following default credentials:

* username: admin
* password: b3secure (this will be soon replaced with random password)

Please change password to your own immediately!

### Prod

Requirements:
* docker and docker-compose
* on firewall open ports: 1194/udp and 8080/tcp

Execute commands

    curl -O https://raw.githubusercontent.com/adamwalach/openvpn-web-ui/master/docs/docker-compose.yml
    docker-compose up -d

It starts two docker containers. One with OpenVPN server and second with OpenVPNAdmin web application. Through a docker volume it creates following directory structure:


    .
    ├── docker-compose.yml
    └── openvpn-data
        ├── conf
        │   ├── dh2048.pem
        │   ├── ipp.txt
        │   ├── keys
        │   │   ├── 01.pem
        │   │   ├── ca.crt
        │   │   ├── ca.key
        │   │   ├── index.txt
        │   │   ├── index.txt.attr
        │   │   ├── index.txt.old
        │   │   ├── serial
        │   │   ├── serial.old
        │   │   ├── server.crt
        │   │   ├── server.csr
        │   │   ├── server.key
        │   │   └── vars
        │   ├── openvpn.log
        │   └── server.conf
        └── db
            └── data.db



### Dev

Requirements:
* golang environments
* [beego](https://beego.me/docs/install/)

Execute commands:

    go get github.com/adamwalach/openvpn-web-ui
    cd $GOPATH/src/github.com/adamwalach/openvpn-web-ui
    bee run -gendoc=true

## Todo

* add unit tests
* add option to modify certificate properties
* generate random admin password at initialization phase
* add versioning
* add automatic ssl/tls (check how [ponzu](https://github.com/ponzu-cms/ponzu) did it)


## License

This project uses [MIT license](LICENSE)

## Remarks

### Vendoring
https://github.com/kardianos/govendor is used for vendoring.

To update dependencies from GOPATH:

`govendor update +v`

### Template
AdminLTE - dashboard & control panel theme. Built on top of Bootstrap 3.

Preview: https://almsaeedstudio.com/themes/AdminLTE/index2.html

