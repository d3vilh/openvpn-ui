// Package routers defines application routes
// @APIVersion 1.0.0
// @Title OpenVPN API
// @Description REST API allows you to control and monitor your OpenVPN server
// @Contact adam.walach@gmail.com
// License Apache 2.0
// LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/controllers"
)

func Init(configDir string) {
	web.SetStaticPath("/swagger", "swagger")
	web.Router("/", &controllers.MainController{})
	web.Router("/login", &controllers.LoginController{}, "get:Login;post:Login")
	web.Router("/logout", &controllers.LoginController{}, "get:Logout")
	web.Router("/auth/google", &controllers.LoginController{}, "get:GoogleLogin")
	web.Router("/auth/google/callback", &controllers.LoginController{}, "get:GoogleCallback")	
	web.Router("/profile", &controllers.ProfileController{})
	web.Router("/settings", &controllers.SettingsController{})
	web.Router("/ov/config", &controllers.OVConfigController{})
	web.Router("/logs", &controllers.LogsController{})
	web.Router("/ov/clientconfig", &controllers.OVClientConfigController{ConfigDir: configDir})
	web.Router("/easyrsa/config", &controllers.EasyRSAConfigController{ConfigDir: configDir})
	web.Router("/dangerzone", &controllers.DangerController{})

	web.Include(&controllers.CertificatesController{ConfigDir: configDir})
	web.Include(&controllers.DangerController{})
	web.Include(&controllers.OVConfigController{ConfigDir: configDir})
	web.Include(&controllers.OVClientConfigController{ConfigDir: configDir})
	web.Include(&controllers.ProfileController{})

	ns := web.NewNamespace("/api/v1",
		web.NSNamespace("/session",
			web.NSInclude(
				&controllers.APISessionController{},
			),
		),
		web.NSNamespace("/sysload",
			web.NSInclude(
				&controllers.APISysloadController{},
			),
		),
		web.NSNamespace("/signal",
			web.NSInclude(
				&controllers.APISignalController{},
			),
		),
	)
	web.AddNamespace(ns)
}
