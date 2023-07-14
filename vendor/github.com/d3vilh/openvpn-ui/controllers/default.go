package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	mi "github.com/d3vilh/openvpn-server-config/server/mi"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/state"
)

type MainController struct {
	BaseController
}

func (c *MainController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Status",
	}
}

func (c *MainController) Get() {
	c.Data["sysinfo"] = lib.GetSystemInfo()
	lib.Dump(lib.GetSystemInfo())
	client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
	status, err := client.GetStatus()
	if err != nil {
		beego.Error(err)
		beego.Warn(fmt.Sprintf("passed client line: %s", client))
		beego.Warn(fmt.Sprintf("error: %s", err))
	} else {
		c.Data["ovstatus"] = status
	}
	lib.Dump(status)

	version, err := client.GetVersion()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovversion"] = version.OpenVPN
	}
	lib.Dump(version)

	pid, err := client.GetPid()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovpid"] = pid
	}
	lib.Dump(pid)

	loadStats, err := client.GetLoadStats()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovstats"] = loadStats
	}
	lib.Dump(loadStats)

	c.TplName = "index.html"
}
