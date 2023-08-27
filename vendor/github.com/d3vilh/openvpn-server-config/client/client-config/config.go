package clientconfig

import (
	"bytes"
	"os"
	"text/template"
)

var defaultConfig = Config{
	Device:            "tun",
	Proto:             "udp",
	ServerAddress:     "127.0.0.1",
	Port:              1194,
	OpenVpnServerPort: "12235",
	Cipher:            "AES-256-CBC",
	Auth:              "SHA256",
	Ca:                "ca.crt",
}

// Config model
type Config struct {
	Device            string
	ServerAddress     string
	Port              int
	OpenVpnServerPort string
	Proto             string

	Ca   string
	Cert string
	Key  string
	Ta   string

	Cipher string
	Auth   string
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
