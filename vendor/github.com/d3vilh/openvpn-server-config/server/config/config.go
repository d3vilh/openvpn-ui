package config

// html/template changed to text/template
import (
	"bytes"
	"io/ioutil"
	"text/template"
)

// Don't think these defaults are ever used -- see models/models.go
var defaultConfig = Config{
	Management:          "0.0.0.0 2080",
	Port:                1194,
	ClientPort:          12235,
	Proto:               "udp",
	Device:              "tun",
	Ca:                  "pki/ca.crt",
	Cert:                "pki/issued/server.crt",
	Key:                 "pki/private/server.key",
	Cipher:              "AES-256-CBC",
	Keysize:             256,
	Auth:                "SHA512",
	Dh:                  "pki/dh.pem",
	Server:              "10.0.70.0 255.255.255.0",
	Route:               "10.0.71.0 255.255.255.0",
	IfconfigPoolPersist: "pki/ipp.txt",
	PushRoute:           "route \"10.0.60.0 255.255.255.0\"",
	DNSServer1:          "dhcp-option DNS 8.8.8.8",
	DNSServer2:          "dhcp-option DNS 1.0.0.1",
	Keepalive:           "10 120",
	MaxClients:          100,
}

// Config model
type Config struct {
	Management string
	Port       int
	ClientPort int
	Proto      string
	Device     string

	Ca   string
	Cert string
	Key  string

	Cipher  string
	Keysize int
	Auth    string
	Dh      string

	Server              string
	Route               string
	IfconfigPoolPersist string
	PushRoute           string
	DNSServer1          string
	DNSServer2          string
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
	template, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(template), c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destPath, []byte(str), 0644)
}
