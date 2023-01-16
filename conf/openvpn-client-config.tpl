client
dev {{ .Device }}
proto {{ .Proto}}
remote {{ .ServerAddress }} {{ .OpenVpnServerPort }} {{ .Proto }}
resolv-retry infinite
user nobody
group nogroup
persist-tun
persist-key
remote-cert-tls server
cipher {{ .Cipher }}
keysize 256
auth {{ .Auth }}
auth-nocache
tls-client
redirect-gateway def1
comp-lzo
verb 3
<ca>
{{ .Ca }}
</ca>
<cert>
{{ .Cert }}
</cert>
<key>
{{ .Key }}
</key>