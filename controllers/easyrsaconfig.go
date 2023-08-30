package controllers

import (
	"html/template"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	easyrsaconfig "github.com/d3vilh/openvpn-server-config/easyrsa/config"
	mi "github.com/d3vilh/openvpn-server-config/server/mi"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/state"
)

type EasyRSAConfigController struct {
	BaseController
	ConfigDir string
}

func (c *EasyRSAConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "EasyRSA configuration",
	}
}

func (c *EasyRSAConfigController) Get() {
	c.TplName = "easyrsavar.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	cfg := models.EasyRSAConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *EasyRSAConfigController) Post() {
	c.TplName = "easyrsavar.html"
	flash := web.NewFlash()
	cfg := models.EasyRSAConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPath := filepath.Join(state.GlobalCfg.EasyRSAPath, "pki/vars")
	err := easyrsaconfig.SaveToFile(filepath.Join(c.ConfigDir, "easyrsa-vars.tpl"), cfg.Config, destPath)
	if err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	destPath = filepath.Join(state.GlobalCfg.OVConfigPath, "config/easy-rsa.vars")
	err = easyrsaconfig.SaveToFile(filepath.Join(c.ConfigDir, "easyrsa-vars.tpl"), cfg.Config, destPath)
	if err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}
	flash.Store(&c.Controller)

}
