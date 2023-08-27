package models

import (
	"github.com/beego/beego/v2/client/orm"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
)

// OVClientConfig holds values for OpenVPN Client config file
type OVClientConfig struct {
	Id      int
	Profile string `orm:"size(64);unique" valid:"Required;"`
	//	ServerAddress     string `orm:"size(64);unique" form:"ServerAddress" valid:"Required;"`
	//	OpenVpnServerPort string `orm:"size(64);unique" form:"OpenVpnServerPort" valid:"Required;"`
	//	Proto             string `orm:"size(64);unique" form:"Proto" valid:"Required;"`
	clientconfig.Config
}

// Insert wrapper
func (c *OVClientConfig) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}

	return nil
}

// Read wrapper
func (c *OVClientConfig) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

// Update wrapper
func (c *OVClientConfig) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

// Delete wrapper
func (c *OVClientConfig) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}
