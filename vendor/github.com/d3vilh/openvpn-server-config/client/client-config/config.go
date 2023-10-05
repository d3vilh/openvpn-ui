package clientconfig

import (
	"bytes"
	"os"
	"text/template"
)

var defaultConfig = Config{
	FuncMode:          0, // 0 = standard authentication (cert, cert + password), 1 = 2FA authentication (cert + OTP)
	Device:            "tun",
	Proto:             "udp",
	ServerAddress:     "127.0.0.1",
	Port:              1194,
	OpenVpnServerPort: "1194",
	Cipher:            "AES-256-GCM",
	RedirectGateway:   "redirect-gateway def1",
	Auth:              "SHA256",
	Ca:                "ca.crt",
	AuthUserPass:      "",                 // "auth-user-pass" when 2fa
	TFAIssuer:         "MFA%20OpenVPN-UI", // 2FA issuer
	CustomConfOne:     "#Custom Option One",
	CustomConfTwo:     "#Custom Option Two",
	CustomConfThree:   "#Custom Option Three",
}

// Config model
type Config struct {
	FuncMode          int
	Device            string
	ServerAddress     string
	Port              int
	OpenVpnServerPort string
	Proto             string

	Ca   string
	Cert string
	Key  string
	Ta   string

	Cipher          string
	RedirectGateway string
	Auth            string
	AuthUserPass    string
	TFAIssuer       string

	CustomConfOne   string
	CustomConfTwo   string
	CustomConfThree string
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
