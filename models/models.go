package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
	easyrsaconfig "github.com/d3vilh/openvpn-server-config/easyrsa/config"
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
		new(EasyRSAConfig),
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
		IsAdmin:  true,
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

	easyRSAPath, err := web.AppConfig.String("EasyRsaPath")
	if err != nil {
		return nil, err
	}

	s := Settings{
		Profile:      "default",
		MIAddress:    miAddress,
		MINetwork:    miNetwork,
		OVConfigPath: ovConfigPath,
		EasyRSAPath:  easyRSAPath,
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
			FuncMode:   0,
			Management: fmt.Sprintf("%s %s", address, network),
			//	ScriptSecurity:           "#! script-security 2",
			//	UserPassVerify:           "#! auth-user-pass-verify /opt/app/bin/oath.sh via-file",
			ScriptSecurity:           "",
			UserPassVerify:           "",
			Device:                   "tun",
			Port:                     1194,
			Proto:                    "udp",
			OVConfigTopology:         "subnet",
			Keepalive:                "10 120",
			MaxClients:               100,
			OVConfigUser:             "nobody",
			OVConfigGroup:            "nogroup",
			OVConfigClientConfigDir:  "/etc/openvpn/staticclients",
			IfconfigPoolPersist:      "pki/ipp.txt",
			Ca:                       "pki/ca.crt",
			Cert:                     "pki/issued/server.crt",
			Key:                      "pki/private/server.key",
			Crl:                      "pki/crl.pem",
			Dh:                       "pki/dh.pem",
			TLSControlChannel:        "tls-crypt pki/ta.key",
			TLSMinVersion:            "tls-version-min 1.2",
			TLSRemoteCert:            "remote-cert-tls client",
			Cipher:                   "AES-256-GCM",
			OVConfigNcpCiphers:       "AES-256-GCM:AES-192-GCM:AES-128-GCM",
			Auth:                     "SHA512",
			Server:                   "server 10.0.70.0 255.255.255.0",
			Route:                    "route 10.0.71.0 255.255.255.0",
			PushRoute:                "push \"route 10.0.60.0 255.255.255.0\"",
			DNSServer1:               "push \"dhcp-option DNS 8.8.8.8\"",
			DNSServer2:               "push \"dhcp-option DNS 1.0.0.1\"",
			RedirectGW:               "push \"redirect-gateway def1 bypass-dhcp\"",
			OVConfigLogfile:          "/var/log/openvpn/openvpn.log",
			OVConfigLogVerbose:       3,
			OVConfigStatusLog:        "/var/log/openvpn/openvpn-status.log",
			OVConfigStatusLogVersion: 2,
			CustomOptOne:             "# Custom Option One",
			CustomOptTwo:             "# Custom Option Two\n# client-to-client",
			CustomOptThree:           "# Custom Option Three\n# push \"route 0.0.0.0 255.255.255.255 net_gateway\"\n# push block-outside-dns",
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
			OpenVpnServerPort: "1194",
			Cipher:            "AES-256-GCM",
			RedirectGateway:   "redirect-gateway def1",
			Auth:              "SHA512",
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

func CreateDefaultEasyRSAConfig(configDir string, easyRSAPath string, address string, network string) {
	c := EasyRSAConfig{
		Profile: "default",
		Config: easyrsaconfig.Config{
			EasyRSADN:          "org",
			EasyRSAReqCountry:  "UA",
			EasyRSAReqProvince: "KY",
			EasyRSAReqCity:     "Kyiv",
			EasyRSAReqOrg:      "SweetHome",
			EasyRSAReqEmail:    "sweet@home.net",
			EasyRSAReqOu:       "MyOrganizationalUnit",
			EasyRSAReqCn:       "server",
			EasyRSAKeySize:     2048,
			EasyRSACaExpire:    3650,
			EasyRSACertExpire:  825,
			EasyRSACertRenew:   30,
			EasyRSACrlDays:     180,
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			logs.Info("New settings profile created")
		} else {
			logs.Debug(c)
		}
		easyRSAConfig := filepath.Join(easyRSAPath, "pki/vars")
		if _, err = os.Stat(easyRSAConfig); os.IsNotExist(err) {
			if err = easyrsaconfig.SaveToFile(filepath.Join(configDir, "easyrsa-vars.tpl"), c.Config, easyRSAConfig); err != nil {
				logs.Error(err)
			}
		}
	} else {
		logs.Error(err)
	}
}
