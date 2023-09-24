package controllers

import (
	"html/template"

	passlib "gopkg.in/hlandau/passlib.v1"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
)

type ProfileController struct {
	BaseController
}

func (c *ProfileController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Profile",
	}
}

func (c *ProfileController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["profile"] = c.Userinfo
	c.TplName = "profile.html"
}

func (c *ProfileController) Post() {
	c.TplName = "profile.html"
	c.Data["profile"] = c.Userinfo

	flash := web.NewFlash()

	user := models.User{}
	if err := c.ParseForm(&user); err != nil {
		logs.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user.Login = c.Userinfo.Login
	c.Data["profile"] = user

	if vMap := validateUser(user); vMap != nil {
		c.Data["validation"] = vMap
		return
	}

	hash, err := passlib.Hash(user.Password)
	if err != nil {
		flash.Error("Unable to hash password")
		flash.Store(&c.Controller)
		return
	}
	c.Userinfo.Email = user.Email
	c.Userinfo.Name = user.Name
	c.Userinfo.Password = hash
	o := orm.NewOrm()
	if _, err := o.Update(c.Userinfo); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Profile has been updated")
	}
	flash.Store(&c.Controller)
}

func validateUser(user models.User) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&user)
	if err != nil {
		logs.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func (c *ProfileController) Create() {
	c.TplName = "profile.html"
	c.Data["profile"] = c.Userinfo

	flash := web.NewFlash()
	user := models.User{}
	if err := c.ParseForm(&user); err != nil {
		logs.Error(err)
		return
	}
	logs.Info("Creating new user with the following information:")
	logs.Info("Login:", user.Login)
	logs.Info("Name:", user.Name)
	logs.Info("Email:", user.Email)
	logs.Info("Password:", user.Password)
	CreateNewUser(user.Login, user.Name, user.Email, user.Password)
	logs.Info("Creating complete. Enjoy!")
	// Redirect to the newly created user's profile page
	// c.Ctx.Redirect(302, "/profile/"+user.Login)
	flash.Store(&c.Controller)
}

// Create new user
func CreateNewUser(NewLogin string, NewName string, NewEmail string, NewPassword string) {
	o := orm.NewOrm()
	var lastUser models.User
	err := o.QueryTable("user").OrderBy("-id").One(&lastUser)
	if err == orm.ErrNoRows {
		lastUser.Id = 0
	} else if err != nil {
		logs.Error(err)
		return
	}
	newUser := models.User{
		Id:       lastUser.Id + 1,
		Login:    NewLogin,
		Name:     NewName,
		Email:    NewEmail,
		Password: NewPassword,
	}
	hash, err := passlib.Hash(newUser.Password)
	if err != nil {
		logs.Error("Unable to hash password", err)
		return
	}
	newUser.Password = hash
	if created, _, err := o.ReadOrCreate(&newUser, "Name"); err == nil {
		if created {
			logs.Info("New user account created")
		} else {
			logs.Debug(newUser)
		}
	} else {
		logs.Error(err)
	}
}
