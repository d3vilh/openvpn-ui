package models

import (
	"github.com/beego/beego/v2/client/orm"
	easyrsaconfig "github.com/d3vilh/openvpn-server-config/easyrsa/config"
)

// OVClientConfig holds values for OpenVPN Client config file
type EasyRSAConfig struct {
	Id      int
	Profile string `orm:"size(64);unique" valid:"Required;"`
	easyrsaconfig.Config
}

// Insert wrapper
func (c *EasyRSAConfig) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}

	return nil
}

// Read wrapper
func (c *EasyRSAConfig) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

// Update wrapper
func (c *EasyRSAConfig) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

// Delete wrapper
func (c *EasyRSAConfig) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}
