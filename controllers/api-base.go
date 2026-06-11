package controllers

import (
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type APIBaseController struct {
	BaseController
}

type APITokenController struct {
	web.Controller
}

func (c *APITokenController) Prepare() {
	c.EnableXSRF = false
	if !c.validateToken() {
		c.serveUnauthorized()
	}
}

func (c *APITokenController) validateToken() bool {
	apiToken := getAPIToken()
	if apiToken == "" {
		logs.Warning("ApiToken is not configured")
		return false
	}

	auth := c.Ctx.Input.Header("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ") == apiToken
	}
	return c.Ctx.Input.Header("X-API-Token") == apiToken
}

func getAPIToken() string {
	if v := os.Getenv("OPENVPN_API_TOKEN"); v != "" {
		return v
	}
	apiToken, err := web.AppConfig.String("ApiToken")
	if err != nil {
		return ""
	}
	return apiToken
}

func (c *APITokenController) serveUnauthorized() {
	c.Ctx.Output.SetStatus(401)
	c.Data["json"] = JSONResponse{
		Status:  "error",
		Message: "Unauthorized",
	}
	c.ServeJSON()
	c.StopRun()
}

func (c *APITokenController) ServeJSONMessage(message string) {
	r := NewJSONResponse()
	r.Message = message
	c.Data["json"] = r
	c.ServeJSON()
}

func (c *APITokenController) ServeJSONData(data interface{}) {
	r := NewJSONResponse()
	r.Data = data
	c.Data["json"] = r
	c.ServeJSON()
}

func (c *APITokenController) ServeJSONError(message string) {
	c.Data["json"] = JSONResponse{
		Status:  "error",
		Message: message,
	}
	logs.Warning(message)
	c.Ctx.Output.SetStatus(400)
	c.ServeJSON()
}

// JSONResponse http://stackoverflow.com/a/12979961
type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewJSONResponse() *JSONResponse {
	response := &JSONResponse{
		Status: "success",
	}
	return response
}

func (c *APIBaseController) Prepare() {
	c.EnableXSRF = false
	c.BaseController.Prepare()
}

func (c *APIBaseController) NestPrepare() {
	if !c.IsLogin {
		c.ServeJSONError("You are not authorized")
		return
	}
}

func (c *APIBaseController) ServeJSONMessage(message string) {
	r := NewJSONResponse()
	r.Message = message
	c.Data["json"] = r
	c.ServeJSON()
}

func (c *APIBaseController) ServeJSONData(data interface{}) {
	r := NewJSONResponse()
	r.Data = data
	c.Data["json"] = r
	c.ServeJSON()
}

func (c *APIBaseController) ServeJSONError(message string) {
	c.Data["json"] = JSONResponse{
		Status:  "error",
		Message: message,
	}
	logs.Warning(message)
	c.Ctx.Output.SetStatus(400)
	c.ServeJSON()
}
