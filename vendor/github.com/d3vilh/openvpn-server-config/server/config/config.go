package config

// html/template changed to text/template
import (
	"bytes"
	"os"
	"text/template"
)

// Don't think these defaults are ever used -- see models/models.go
var defaultConfig = Config{
	Management:          "0.0.0.0 2080",
	Port:                1194,
	Proto:               "udp",
	Device:              "tun",
	Ca:                  "pki/ca.crt",
	Cert:                "pki/issued/server.crt",
	Key:                 "pki/private/server.key",
	Cipher:              "AES-256-CBC",
	Auth:                "SHA512",
	Dh:                  "pki/dh.pem",
	Server:              "10.0.70.0 255.255.255.0",
	Route:               "10.0.71.0 255.255.255.0",
	IfconfigPoolPersist: "pki/ipp.txt",
	OVConfigLogV:        3,
	OVConfigLogVersion:  2,
	PushRoute:           "10.0.60.0 255.255.255.0",
	DNSServer1:          "8.8.8.8",
	DNSServer2:          "1.0.0.1",
	RedirectGW:          "push \"redirect-gateway def1 bypass-dhcp\"",
	Keepalive:           "10 120",
	MaxClients:          100,
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

	Server              string
	Route               string
	IfconfigPoolPersist string
	OVConfigLogV        int
	OVConfigLogVersion  int
	PushRoute           string
	DNSServer1          string
	DNSServer2          string
	RedirectGW          string
	Keepalive           string
	MaxClients          int
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
