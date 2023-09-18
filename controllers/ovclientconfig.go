package controllers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
	mi "github.com/d3vilh/openvpn-server-config/server/mi"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/state"
)

type OVClientConfigController struct {
	BaseController
	ConfigDir string
}

func (c *OVClientConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "OpenVPN Client configuration",
	}
}

func (c *OVClientConfigController) Get() {
	c.TplName = "ovclient.html"

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "config/client.conf")
	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	cfg := models.OVClientConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *OVClientConfigController) Post() {
	c.TplName = "ovclient.html"
	flash := web.NewFlash()
	cfg := models.OVClientConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "config/client.conf")
	err := clientconfig.SaveToFile(filepath.Join(c.ConfigDir, "openvpn-client-config.tpl"), cfg.Config, destPath)
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

	destPathClientTempl := filepath.Join(state.GlobalCfg.OVConfigPath, "config/client.conf")
	clientTemplate, err := os.ReadFile(destPathClientTempl)
	if err != nil {
		logs.Error(err)
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Store(&c.Controller)
}
