package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/routers"
	"github.com/d3vilh/openvpn-ui/state"
)

func main() {
	configDir := flag.String("config", "conf", "Path to config dir")
	flag.Parse()

	configFile := filepath.Join(*configDir, "app.conf")
	fmt.Println("Config file:", configFile)

	if err := web.LoadAppConfig("ini", configFile); err != nil {
		panic(err)
	}

	models.InitDB()
	models.CreateDefaultUsers()
	defaultSettings, err := models.CreateDefaultSettings()
	if err != nil {
		panic(err)
	}

	models.CreateDefaultOVConfig(*configDir, defaultSettings.OVConfigPath, defaultSettings.MIAddress, defaultSettings.MINetwork)
	models.CreateDefaultOVClientConfig(*configDir, defaultSettings.OVConfigPath, defaultSettings.MIAddress, defaultSettings.MINetwork)
	state.GlobalCfg = *defaultSettings

	routers.Init(*configDir)

	lib.AddFuncMaps()
	web.Run()
}
