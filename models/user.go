package models

import (
	"time"

	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/core/validation"

	//Sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id            int64
	Login         string    `orm:"size(64);unique" form:"Login" valid:"Required;"`
	Name          string    `orm:"size(64);unique" form:"Name" valid:"Required;"`
	Email         string    `orm:"size(64);unique" form:"Email" valid:"Required;Email"`
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MinSize(6)"`
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required"`
	Lastlogintime time.Time `orm:"type(datetime);null" form:"-"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
}

func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		_ = v.SetError("Repassword", "Passwords do not match")
	}
}

func (u *User) Insert() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

// Read wrapper
func (u *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(u, fields...); err != nil {
		return err
	}
	return nil
}

// Update wrapper
func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

// Delete wrapper
func (u *User) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}
