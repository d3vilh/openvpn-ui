package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type APIBaseController struct {
	BaseController
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
