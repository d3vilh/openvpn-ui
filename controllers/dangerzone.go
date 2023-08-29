package controllers

import (
	"github.com/d3vilh/openvpn-ui/lib"
)

// @router /pki/delete [DeletePKI]
func (c *CertificatesController) DeletePKI() {
	lib.DeletePKI()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	// return
}

// @router /pki/init [InitPKI]
func (c *CertificatesController) InitPKI() {
	lib.InitPKI()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	// return
}
