package controllers

import (
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

// @router /pki/delete [DeletePKI]
func (c *DangerController) DeletePKI() {
	lib.DeletePKI()
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}

// @router /pki/init [InitPKI]
func (c *DangerController) InitPKI() {
	lib.InitPKI()
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}

// @router /container/restart [RestartContainer]
func (c *DangerController) RestartContainer() {
	lib.RestartContainer()
	c.Redirect(c.URLFor("DangerController.Get"), 302)
	// return
}
