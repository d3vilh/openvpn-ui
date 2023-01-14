package controllers

import (
	"encoding/json"

	mi "github.com/d3vilh/openvpn-server-config/server/mi"
	"github.com/d3vilh/openvpn-ui/state"
)

// APISignalController sends signals to OpenVPN daemon
type APISignalController struct {
	APIBaseController
}

// KillParams contains CommonName of session to kill
type SignalParams struct {
	Sname string `json:"sname"`
}

// Send signal to OpenVPN daemon
// @Title Send signal
// @Description Sends signal to OpenVPN daemon
// @Param    body     body     controllers.SignalParams     true      "Signal to send"
// @Success 200 request success
// @Failure 400 request failure
// @router / [post]
func (c *APISignalController) Send() {
	client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
	p := SignalParams{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	if err := client.Signal(p.Sname); err != nil {
		c.ServeJSONError(err.Error())
		return
	}

	c.ServeJSONMessage("Signal sent")
}
