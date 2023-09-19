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
	besettings := models.Settings{Profile: "default"}
	_ = besettings.Read("Profile")
	c.Data["BeeSettings"] = &besettings

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
	//logs.Info("Starting Post method in OVClientConfigController")

	c.TplName = "ovclient.html"
	flash := web.NewFlash()
	cfg := models.OVClientConfig{Profile: "default"}
	_ = cfg.Read("Profile")

	//logs.Info("Post: Parsing form data")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	//logs.Info("Post: Dumping configuration data")
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "config/client.conf")
	//logs.Info("Post: Saving configuration to file according to template")
	err := clientconfig.SaveToFile(filepath.Join(c.ConfigDir, "openvpn-client-config.tpl"), cfg.Config, destPath)
	if err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	//logs.Info("Post: Updating configuration in database")
	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Post: Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}

	//logs.Info("Post: Reading updated server configuration from file")
	clientTemplate, err := os.ReadFile(destPath)
	if err != nil {
		logs.Error("Error reading Client template from file:", err)
		flash.Error("Error reading Client template from file")
		return
	}
	c.Data["ClientTemplate"] = string(clientTemplate)

	flash.Store(&c.Controller)
}

// @router /ov/clientconfig/edit [Edit]
func (c *OVClientConfigController) Edit() {
	c.TplName = "ovclient.html"
	flash := web.NewFlash()
	cfg := models.OVClientConfig{Profile: "default"}
	_ = cfg.Read("Profile")

	//logs.Info("Post: Parsing form data")
	if err := c.ParseForm(&cfg); err != nil {
		logs.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	//logs.Info("Post: Dumping configuration data")
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	//logs.Info("Starting Edit method in OVClientConfigController")
	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "config/client.conf")

	err := lib.ConfSaveToFile(destPath, c.GetString("ClientTemplate"))
	if err != nil {
		logs.Error("Error saving Client template to file:", err)
		flash.Error("Error saving Client template to file")
		return
	} else {
		//logs.Info("Edit: Client template saved to file:", destPath)
		flash.Success("Client template has been updated")
	}

	clientTempl, err := os.ReadFile(destPath)
	if err != nil {
		logs.Error("Error reading Client template from file:", err)
		flash.Error("Error reading Client template from file")
		return
	}
	c.Data["ClientTemplate"] = string(clientTempl)

	flash.Store(&c.Controller)
}
