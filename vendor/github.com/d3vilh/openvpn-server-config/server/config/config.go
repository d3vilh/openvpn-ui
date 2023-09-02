package config

// html/template changed to text/template
import (
	"bytes"
	"os"
	"text/template"
)

// Don't think these defaults are ever used -- see models/models.go
var defaultConfig = Config{
	Management:               "0.0.0.0 2080",
	Device:                   "tun",
	Port:                     1194,
	Proto:                    "udp",
	OVConfigTopology:         "subnet",
	Keepalive:                "10 120",
	MaxClients:               100,
	OVConfigUser:             "nobody",
	OVConfigGroup:            "nogroup",
	OVConfigClientConfigDir:  "/etc/openvpn/staticclients",
	IfconfigPoolPersist:      "pki/ipp.txt",
	Ca:                       "pki/ca.crt",
	Cert:                     "pki/issued/server.crt",
	Key:                      "pki/private/server.key",
	Crl:                      "pki/crl.pem",
	Dh:                       "pki/dh.pem",
	ECDhCurve:                "ecdh-curve prime256v1",
	TLSControlChannel:        "tls-crypt pki/ta.key",
	TLSMinVersion:            "tls-version-min 1.2",
	TLSRemoteCert:            "remote-cert-tls client",
	Cipher:                   "AES-256-CBC",
	OVConfigNcpCiphers:       "AES-256-GCM:AES-192-GCM:AES-128-GCM",
	Auth:                     "SHA512",
	Server:                   "server 10.0.70.0 255.255.255.0",
	Route:                    "route 10.0.71.0 255.255.255.0",
	PushRoute:                "push \"route 10.0.60.0 255.255.255.0\"",
	DNSServer1:               "push \"dhcp-option DNS 8.8.8.8\"",
	DNSServer2:               "push \"dhcp-option DNS 1.0.0.1\"",
	RedirectGW:               "push \"redirect-gateway def1 bypass-dhcp\"",
	OVConfigLogfile:          "/var/log/openvpn/openvpn.log",
	OVConfigLogVerbose:       3,
	OVConfigStatusLog:        "/var/log/openvpn/openvpn-status.log",
	OVConfigStatusLogVersion: 2,
	CustomOptOne:             "#Custom Option One",
	CustomOptTwo:             "#Custom Option Two",
	CustomOptThree:           "#Custom Option Three",
}

// Config model
type Config struct {
	Management string
	Device     string
	Port       int
	Proto      string

	OVConfigTopology string
	Keepalive        string
	MaxClients       int

	OVConfigUser  string
	OVConfigGroup string

	OVConfigClientConfigDir string
	IfconfigPoolPersist     string

	Ca        string
	Cert      string
	Key       string
	Crl       string
	Dh        string
	ECDhCurve string

	TLSControlChannel string
	TLSMinVersion     string
	TLSRemoteCert     string

	Cipher             string
	OVConfigNcpCiphers string

	Auth string

	Server     string
	Route      string
	PushRoute  string
	DNSServer1 string
	DNSServer2 string
	RedirectGW string

	OVConfigLogfile          string
	OVConfigLogVerbose       int
	OVConfigStatusLog        string
	OVConfigStatusLogVersion int

	CustomOptOne   string
	CustomOptTwo   string
	CustomOptThree string
}

// New returns config object with default values
func New() Config {
	return defaultConfig
}

// GetText injects config values into template
func GetText(tpl string, c Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, c)
	return buf.String(), nil
}

// SaveToFile reads teamplate and writes result to destination file
func SaveToFile(tplPath string, c Config, destPath string) error {
	template, err := os.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(template), c)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, []byte(str), 0644)
}
