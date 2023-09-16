package controllers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/state"
)

type NewCertParams struct {
	Name       string `form:"Name" valid:"Required;"`
	Staticip   string `form:"staticip"`
	Passphrase string `form:"passphrase"`
	ExpireDays string `form:"EasyRSACertExpire"`
	Email      string `form:"EasyRSAReqEmail"`
	Country    string `form:"EasyRSAReqCountry"`
	Province   string `form:"EasyRSAReqProvince"`
	City       string `form:"EasyRSAReqCity"`
	Org        string `form:"EasyRSAReqOrg"`
	OrgUnit    string `form:"EasyRSAReqOu"`
}

type CertificatesController struct {
	BaseController
	ConfigDir string
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Certificates",
	}
}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	keysPath := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/issued")

	cfgPath, err := c.saveClientConfig(keysPath, name)
	if err != nil {
		logs.Error(err)
		return
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		logs.Error(err)
		return
	}
	if _, err = c.Controller.Ctx.ResponseWriter.Write(data); err != nil {
		logs.Error(err)
	}
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
	cfg := models.EasyRSAConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["EasyRSA"] = &cfg
}

func (c *CertificatesController) showCerts() {
	path := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/index.txt")
	certs, err := lib.ReadCerts(path)
	if err != nil {
		logs.Error(err)
	}
	lib.Dump(certs)
	c.Data["certificates"] = &certs
}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.CreateCertificate(cParams.Name, cParams.Staticip, cParams.Passphrase, cParams.ExpireDays, cParams.Email, cParams.Country, cParams.Province, cParams.City, cParams.Org, cParams.OrgUnit); err != nil {
				logs.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			} else {
				flash.Success("Success! Certificate for the name \"" + cParams.Name + "\" has been created")
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Revoke() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	serial := c.GetString(":serial")
	if err := lib.RevokeCertificate(name, serial); err != nil {
		logs.Error(err)
		//flash.Error(err.Error())
		//flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Certificate for the name \"" + name + "\" and serial  \"" + serial + "\" has been revoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/restart [get]
func (c *CertificatesController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	// return
}

// @router /certificates/burn/:key/:serial [get]
func (c *CertificatesController) Burn() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	CN := c.GetString(":key")
	serial := c.GetString(":serial")
	if err := lib.BurnCertificate(CN, serial); err != nil {
		logs.Error(err)
		//flash.Error(err.Error())
		//flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Certificate for the name \"" + CN + "\" and serial  \"" + serial + "\"  has been removed")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Renew() {
	c.TplName = "certificates.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	localip := c.GetString(":localip")
	serial := c.GetString(":serial")
	if err := lib.RenewCertificate(name, localip, serial); err != nil {
		logs.Error(err)
		//flash.Error(err.Error())
		//flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Certificate for the name \"" + name + "\"  and \"" + localip + "\" and \"" + serial + "\" has been renewed")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func (c *CertificatesController) saveClientConfig(keysPath string, name string) (string, error) {
	cfg := clientconfig.New()
	keysPathCa := filepath.Join(state.GlobalCfg.OVConfigPath, "pki")
	ServerAddress := models.OVClientConfig{Profile: "default"}
	_ = ServerAddress.Read("Profile")
	cfg.ServerAddress = ServerAddress.ServerAddress
	OpenVpnServerPort := models.OVClientConfig{Profile: "default"}
	_ = OpenVpnServerPort.Read("Profile")
	cfg.OpenVpnServerPort = OpenVpnServerPort.OpenVpnServerPort

	ca, err := os.ReadFile(filepath.Join(keysPathCa, "ca.crt"))
	if err != nil {
		return "", err
	}
	cfg.Ca = string(ca)

	ta, err := os.ReadFile(filepath.Join(keysPathCa, "ta.key"))
	if err != nil {
		return "", err
	}
	cfg.Ta = string(ta)

	cert, err := os.ReadFile(filepath.Join(keysPath, name+".crt"))
	if err != nil {
		return "", err
	}
	cfg.Cert = string(cert)

	keysPathKey := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/private")
	key, err := os.ReadFile(filepath.Join(keysPathKey, name+".key"))
	if err != nil {
		return "", err
	}
	cfg.Key = string(key)

	serverConfig := models.OVConfig{Profile: "default"}
	_ = serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "clients", name+".ovpn")
	if err := SaveToFile(filepath.Join(c.ConfigDir, "openvpn-client-config.tpl"), cfg, destPath); err != nil {
		logs.Error(err)
		return "", err
	}

	return destPath, nil
}

func GetText(tpl string, c clientconfig.Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SaveToFile(tplPath string, c clientconfig.Config, destPath string) error {
	tpl, err := os.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(tpl), c)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, []byte(str), 0644)
}
