package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
	"github.com/d3vilh/openvpn-server-config/server/config"
	"gopkg.in/hlandau/passlib.v1"
)

func InitDB() {
	err := orm.RegisterDriver("sqlite3", orm.DRSqlite)
	if err != nil {
		panic(err)
	}
	dbPath, err := web.AppConfig.String("dbPath")
	if err != nil {
		panic(err)
	}
	dbSource := "file:" + dbPath

	err = orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
		new(OVClientConfig),
	)

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error(err)
		return
	}
}

func CreateDefaultUsers() {
	hash, err := passlib.Hash(os.Getenv("OPENVPN_ADMIN_PASSWORD"))
	if err != nil {
		logs.Error("Unable to hash password", err)
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
			logs.Info("Default admin account created")
		} else {
			logs.Debug(user)
		}
	}

}

func CreateDefaultSettings() (*Settings, error) {
	miAddress, err := web.AppConfig.String("OpenVpnManagementAddress")
	if err != nil {
		return nil, err
	}
	miNetwork, err := web.AppConfig.String("OpenVpnManagementNetwork")
	if err != nil {
		return nil, err
	}
	ovConfigPath, err := web.AppConfig.String("OpenVpnPath")
	if err != nil {
		return nil, err
	}

	s := Settings{
		Profile:      "default",
		MIAddress:    miAddress,
		MINetwork:    miNetwork,
		OVConfigPath: ovConfigPath,
		//	ServerAddress:     serverAddress,
		//	OpenVpnServerPort: serverPort,
	}

	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		if created {
			logs.Info("New settings profile created")
		} else {
			logs.Debug(s)
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
			OVConfigLogVersion:  2,
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
			logs.Info("New settings profile created")
		} else {
			logs.Debug(c)
		}
		serverConfig := filepath.Join(ovConfigPath, "config/server.conf")
		if _, err = os.Stat(serverConfig); os.IsNotExist(err) {
			if err = config.SaveToFile(filepath.Join(configDir, "openvpn-server-config.tpl"), c.Config, serverConfig); err != nil {
				logs.Error(err)
			}
		}
	} else {
		logs.Error(err)
	}
}

func CreateDefaultOVClientConfig(configDir string, ovConfigPath string, address string, network string) {
	c := OVClientConfig{
		Profile: "default",
		Config: clientconfig.Config{
			Device:            "tun",
			Port:              1194,
			Proto:             "udp",
			ServerAddress:     "127.0.0.1",
			OpenVpnServerPort: "12235",
			Cipher:            "AES-256-CBC",
			Auth:              "SHA512",
			Ca:                "pki/ca.crt",
			Cert:              "pki/issued/server.crt",
			Key:               "pki/private/server.key",
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			logs.Info("New settings profile created")
		} else {
			logs.Debug(c)
		}
		clientConfig := filepath.Join(ovConfigPath, "config/client.conf")
		if _, err = os.Stat(clientConfig); os.IsNotExist(err) {
			if err = clientconfig.SaveToFile(filepath.Join(configDir, "openvpn-client-config.tpl"), c.Config, clientConfig); err != nil {
				logs.Error(err)
			}
		}
	} else {
		logs.Error(err)
	}
}
