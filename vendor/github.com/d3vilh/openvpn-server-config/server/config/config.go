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
	Port:                     1194,
	Proto:                    "udp",
	Device:                   "tun",
	Ca:                       "pki/ca.crt",
	Cert:                     "pki/issued/server.crt",
	Key:                      "pki/private/server.key",
	Cipher:                   "AES-256-CBC",
	Auth:                     "SHA512",
	Dh:                       "pki/dh.pem",
	Crl:                      "pki/crl.pem",
	Server:                   "10.0.70.0 255.255.255.0",
	Route:                    "10.0.71.0 255.255.255.0",
	IfconfigPoolPersist:      "pki/ipp.txt",
	PushRoute:                "10.0.60.0 255.255.255.0",
	DNSServer1:               "8.8.8.8",
	DNSServer2:               "1.0.0.1",
	RedirectGW:               "push \"redirect-gateway def1 bypass-dhcp\"",
	Keepalive:                "10 120",
	MaxClients:               100,
	OVConfigLogfile:          "/var/log/openvpn/openvpn.log",
	OVConfigLogVerbose:       3,
	OVConfigTopology:         "subnet",
	OVConfigClientConfigDir:  "/etc/openvpn/staticclients",
	OVConfigNcpCiphers:       "AES-256-GCM:AES-192-GCM:AES-128-GCM",
	OVConfigStatusLog:        "/var/log/openvpn/openvpn-status.log",
	OVConfigStatusLogVersion: 2,
	OVConfigUser:             "nobody",
	OVConfigGroup:            "nogroup",
	CustomOptOne:             "#Custom Option One",
	CustomOptTwo:             "#Custom Option Two",
	CustomOptThree:           "#Custom Option Three",
}

// Config model
type Config struct {
	Management string
	Port       int
	Proto      string
	Device     string

	Ca   string
	Cert string
	Key  string

	Cipher string
	Auth   string
	Dh     string
	Crl    string

	Server              string
	Route               string
	IfconfigPoolPersist string
	PushRoute           string
	DNSServer1          string
	DNSServer2          string

	RedirectGW string

	Keepalive                string
	MaxClients               int
	OVConfigLogfile          string
	OVConfigLogVerbose       int
	OVConfigTopology         string
	OVConfigClientConfigDir  string
	OVConfigNcpCiphers       string
	OVConfigStatusLog        string
	OVConfigStatusLogVersion int
	OVConfigUser             string
	OVConfigGroup            string
	CustomOptOne             string
	CustomOptTwo             string
	CustomOptThree           string
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
