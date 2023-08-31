package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/lib"
)

type DangerController struct {
	BaseController
}

func (c *DangerController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
}

func (c *DangerController) Get() {
	c.TplName = "dangerzone.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Danger Zone!",
	}
}

// @router /pki/delete:key [DeletePKI]
func (c *DangerController) DeletePKI() {
	c.TplName = "dangerzone.html"
	flash := web.NewFlash()
	name := c.GetString(":key")
	if err := lib.DeletePKI(name); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Warning("Success! The \"" + name + "\" has been deleted")
		flash.Store(&c.Controller)
	}
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}

// @router /pki/init [InitPKI]
func (c *DangerController) InitPKI() {
	c.TplName = "dangerzone.html"
	flash := web.NewFlash()
	if err := lib.InitPKI(); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Success("Success! PKI has been initialized.")
		flash.Store(&c.Controller)
	}
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}

// @router /container/restart [RestartContainer]
func (c *DangerController) RestartContainer() {
	c.TplName = "dangerzone.html"
	flash := web.NewFlash()
	if err := lib.RestartContainer(); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Container has been restarted")
		flash.Store(&c.Controller)
	}
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}
