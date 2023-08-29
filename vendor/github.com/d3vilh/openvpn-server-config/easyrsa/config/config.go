package easyrsaconfig

// html/template changed to text/template
import (
	"bytes"
	"os"
	"text/template"
)

// Don't think these defaults are ever used -- see models/models.go
var defaultConfig = Config{
	EasyRSADN:          "org",
	EasyRSAReqCountry:  "UA",
	EasyRSAReqProvince: "KY",
	EasyRSAReqCity:     "Kyiv",
	EasyRSAReqOrg:      "SweetHome",
	EasyRSAReqEmail:    "sweet@home.net",
	EasyRSAReqOu:       "MyOrganizationalUnit",
	EasyRSAReqCn:       "server",
	EasyRSAKeySize:     2048,
	EasyRSACaExpire:    3650,
	EasyRSACertExpire:  825,
	EasyRSACertRenew:   30,
	EasyRSACrlDays:     180,
}

// Config model
type Config struct {
	EasyRSADN          string
	EasyRSAReqCountry  string
	EasyRSAReqProvince string
	EasyRSAReqCity     string
	EasyRSAReqOrg      string
	EasyRSAReqEmail    string
	EasyRSAReqOu       string
	EasyRSAReqCn       string

	EasyRSAKeySize    int
	EasyRSACaExpire   int
	EasyRSACertExpire int
	EasyRSACertRenew  int
	EasyRSACrlDays    int
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
