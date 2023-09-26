package routers

import (
	"github.com/beego/beego/v2/server/web"
)

func init() {

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISessionController"],
			web.ControllerComments{
				Method:           "Kill",
				Router:           `/`,
				AllowHTTPMethods: []string{"delete"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISignalController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISignalController"],
			web.ControllerComments{
				Method:           "Send",
				Router:           `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISysloadController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:APISysloadController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Download",
				Router:           `/certificates/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           `/certificates`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Revoke",
				Router:           `/certificates/revoke/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Restart",
				Router:           `/certificates/restart`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Burn",
				Router:           `/certificates/burn/:key/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:CertificatesController"],
			web.ControllerComments{
				Method:           "Renew",
				Router:           `/certificates/renew/:key/:localip/:serial`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"],
			web.ControllerComments{
				Method:           "DeletePKI",
				Router:           `/pki/delete/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"],
			web.ControllerComments{
				Method:           "InitPKI",
				Router:           `/pki/init/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:DangerController"],
			web.ControllerComments{
				Method:           "RestartContainer",
				Router:           `/container/restart/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"],
			web.ControllerComments{
				Method:           "Edit",
				Router:           `/ov/config/edit`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           `/ov/config`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVConfigController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/ov/config`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"],
			web.ControllerComments{
				Method:           "Edit",
				Router:           `/ov/clientconfig/edit`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"],
			web.ControllerComments{
				Method:           "Post",
				Router:           `/ov/clientconfig`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:OVClientConfigController"],
			web.ControllerComments{
				Method:           "Get",
				Router:           `/ov/clientconfig`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:ProfileController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:ProfileController"],
			web.ControllerComments{
				Method:           "Create",
				Router:           `/profile/create`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

	web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:ProfileController"] =
		append(web.GlobalControllerRouter["github.com/d3vilh/openvpn-ui/controllers:ProfileController"],
			web.ControllerComments{
				Method:           "DeleteUser",
				Router:           `/profile/delete/:key`,
				AllowHTTPMethods: []string{"get"},
				Params:           nil})

}
