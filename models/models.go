package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/d3vilh/openvpn-server-config/server/config"
	"gopkg.in/hlandau/passlib.v1"
)

func InitDB() {
	err := orm.RegisterDriver("sqlite3", orm.DRSqlite)
	if err != nil {
		panic(err)
	}
	dbSource := "file:" + beego.AppConfig.String("dbPath")

	err = orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
	)

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Error(err)
		return
	}
}

func CreateDefaultUsers() {
	hash, err := passlib.Hash(os.Getenv("OPENVPN_ADMIN_PASSWORD"))
	if err != nil {
		beego.Error("Unable to hash password", err)
	}
	user := User{
		Id:       1,
		Login:    os.Getenv("OPENVPN_ADMIN_USERNAME"),
		Name:     "Administrator",
		Email:    "root@localhost",
		Password: hash,
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			beego.Info("Default admin account created")
		} else {
			beego.Debug(user)
		}
	}

}

func CreateDefaultSettings() (*Settings, error) {
	s := Settings{
		Profile:   "default",
		MIAddress: beego.AppConfig.String("OpenVpnManagementAddress"),
		//MIAddress:         "openvpn:2080",
		MINetwork:         beego.AppConfig.String("OpenVpnManagementNetwork"),
		ServerAddress:     beego.AppConfig.String("OpenVpnServerAddress"),
		OpenVpnServerPort: beego.AppConfig.String("OpenVpnServerPort"),
		OVConfigPath:      beego.AppConfig.String("OpenVpnPath"),
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(s)
		}
		return &s, nil
	} else {
		return nil, err
	}
}

func CreateDefaultOVConfig(configDir string, ovConfigPath string, address string, network string) {
	c := OVConfig{
		Profile: "default",
		Config: config.Config{
			Device:              "tun",
			Port:                1194,
			Proto:               "udp",
			DNSServer1:          "8.8.8.8",
			DNSServer2:          "1.0.0.1",
			RedirectGW:          "push \"redirect-gateway def1 bypass-dhcp\"",
			PushRoute:           "10.0.60.0 255.255.255.0",
			Route:               "10.0.71.0 255.255.255.0",
			Cipher:              "AES-256-CBC",
			Auth:                "SHA512",
			Dh:                  "pki/dh.pem",
			Keepalive:           "10 120",
			IfconfigPoolPersist: "pki/ipp.txt",
			OVConfigLogV:        3,
			Management:          fmt.Sprintf("%s %s", address, network),
			MaxClients:          100,
			Server:              "10.0.70.0 255.255.255.0",
			Ca:                  "pki/ca.crt",
			Cert:                "pki/issued/server.crt",
			Key:                 "pki/private/server.key",
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(c)
		}
		serverConfig := filepath.Join(ovConfigPath, "config/server.conf")
		if _, err = os.Stat(serverConfig); os.IsNotExist(err) {
			if err = config.SaveToFile(filepath.Join(configDir, "openvpn-server-config.tpl"), c.Config, serverConfig); err != nil {
				beego.Error(err)
			}
		}
	} else {
		beego.Error(err)
	}
}
