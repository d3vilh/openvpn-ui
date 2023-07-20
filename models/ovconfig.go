package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/d3vilh/openvpn-server-config/server/config"
)

// OVConfig holds values for OpenVPN config file
type OVConfig struct {
	Id      int
	Profile string `orm:"size(64);unique" valid:"Required;"`
	config.Config
}

// Insert wrapper
func (c *OVConfig) Insert() error {
	if _, err := orm.NewOrm().Insert(c); err != nil {
		return err
	}

	return nil
}

// Read wrapper
func (c *OVConfig) Read(fields ...string) error {
	if err := orm.NewOrm().Read(c, fields...); err != nil {
		return err
	}
	return nil
}

// Update wrapper
func (c *OVConfig) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

// Delete wrapper
func (c *OVConfig) Delete() error {
	if _, err := orm.NewOrm().Delete(c); err != nil {
		return err
	}
	return nil
}
