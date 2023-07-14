package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"],
			beego.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"],
			beego.ControllerComments{
				Method:           "Kill",
				Router:           `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISignalController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISignalController"],
			beego.ControllerComments{
				Method:           "Send",
				Router:           `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISysloadController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISysloadController"],
			beego.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Download",
				Router:           `/certificates/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Get",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Post",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Revoke",
				Router:           `/certificates/revoke/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Restart",
				Router:           `/certificates/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(beego.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			beego.ControllerComments{
				Method:           "Burn",
				Router:           `/certificates/burn/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})
}
