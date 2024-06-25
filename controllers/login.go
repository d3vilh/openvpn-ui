package controllers

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
)

var (
	oauthConf      *oauth2.Config
	oauthStateString = "random" // use a random string for security purposes
)

func init() {
	clientID, _ := web.AppConfig.String("googleClientID")
	clientSecret, _ := web.AppConfig.String("googleClientSecret")
	redirectURL, _ := web.AppConfig.String("googleRedirectURL")

	oauthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("MainController.Get"))
		return
	}

	c.TplName = "login.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	if !c.Ctx.Input.IsPost() {
		return
	}

	flash := web.NewFlash()
	login := c.GetString("login")
	password := c.GetString("password")

	authType, err := web.AppConfig.String("AuthType")
	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user, err := lib.Authenticate(login, password, authType)

	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	user.Lastlogintime = time.Now()
	err = user.Update("Lastlogintime")
	if err != nil {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	flash.Success("Successfully logged in")
	flash.Store(&c.Controller)

	c.SetLogin(user)

	c.Redirect(c.URLFor("MainController.Get"), 303)
}

func (c *LoginController) Logout() {
	c.DelLogin()
	flash := web.NewFlash()
	flash.Success("Successfully logged out")
	flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("LoginController.Login"))
}

func (c *LoginController) GoogleLogin() {
	url := oauthConf.AuthCodeURL(oauthStateString)
	c.Redirect(url, 302)
}

func (c *LoginController) GoogleCallback() {
	state := c.GetString("state")
	if state != oauthStateString {
		c.Ctx.WriteString("Invalid OAuth state")
		return
	}

	code := c.GetString("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		c.Ctx.WriteString("Code exchange failed: " + err.Error())
		return
	}

	client := oauthConf.Client(context.Background(), token)
	service, err := oauth2api.New(client)
	if err != nil {
		c.Ctx.WriteString("Failed to create OAuth2 service: " + err.Error())
		return
	}

	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.Ctx.WriteString("Failed to get user info: " + err.Error())
		return
	}

	logs.Info("User Info: %+v", userinfo)

	// Get allowed domains from config
	allowedDomainsStr, _ := web.AppConfig.String("allowedDomains")
	allowedDomains := strings.Split(allowedDomainsStr, ",")

	// Check if the user's email domain is allowed
	emailDomain := strings.Split(userinfo.Email, "@")[1]
	allowed := false
	for _, domain := range allowedDomains {
		if emailDomain == domain {
			allowed = true
			break
		}
	}

	if !allowed {
		c.Data["error"] = "Your Email is not allowed to login"
		c.TplName = "login.html"
		c.Render()
		return
	}

	user, err := lib.GetUserByEmail(userinfo.Email)
	if err != nil {
		if err.Error() == "user not found" {
			// Create a new user if not found and set the default values
			user = &models.User{
				Email:         userinfo.Email,
				Name:          userinfo.Email, // Set the name to the email address
				Login:         userinfo.Email,
				Lastlogintime: time.Now(),
				Allowed:       true, // Set to true because authenticated with Google
			}
			err = user.Insert()
			if err != nil {
				c.Ctx.WriteString("Failed to create new user: " + err.Error())
				return
			}
		} else {
			c.Ctx.WriteString("Error fetching user: " + err.Error())
			return
		}
	} else {
		// Update existing user's allowed status, last login time, and name
		user.Allowed = true
		user.Lastlogintime = time.Now()
		user.Name = userinfo.Email // Set the name to the email address
		err = user.Update("Allowed", "Lastlogintime", "Name")
		if err != nil {
			c.Ctx.WriteString("Failed to update user: " + err.Error())
			return
		}
	}

	// Check if the user is allowed
	if !user.Allowed {
		c.Data["error"] = "Access denied"
		c.TplName = "login.html"
		c.Render()
		return
	}

	c.SetLogin(user)

	flash := web.NewFlash()
	flash.Success("Successfully logged in with Google")
	flash.Store(&c.Controller)

	c.Redirect(c.URLFor("MainController.Get"), 302)
}
