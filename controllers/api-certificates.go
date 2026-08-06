package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	clientconfig "github.com/d3vilh/openvpn-server-config/client/client-config"
	"github.com/d3vilh/openvpn-ui/lib"
	"github.com/d3vilh/openvpn-ui/models"
	"github.com/d3vilh/openvpn-ui/state"
)

// APICertificatesController provides token-authenticated certificate management API.
type APICertificatesController struct {
	APITokenController
	ConfigDir string
}

// APICreateCertRequest holds parameters for certificate creation.
type APICreateCertRequest struct {
	Name       string `json:"name" valid:"Required;"`
	StaticIP   string `json:"staticip"`
	Passphrase string `json:"passphrase"`
	ExpireDays string `json:"expire_days"`
	Email      string `json:"email"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Org        string `json:"org"`
	OrgUnit    string `json:"org_unit"`
	TFAName    string `json:"tfa_name"`
	TFAIssuer  string `json:"tfa_issuer"`
}

// Get downloads the .ovpn client configuration for a certificate.
// @Title Get certificate
// @Description Download OpenVPN client configuration (.ovpn)
// @Param    name     path     string     true      "Certificate name"
// @Success 200 file download
// @Failure 400 request failure
// @Failure 401 unauthorized
// @Failure 404 certificate not found
// @router /:name [get]
func (c *APICertificatesController) Get() {
	name := c.GetString(":name")
	if name == "" {
		c.ServeJSONError("name is required")
		return
	}

	cfgPath, err := apiSaveClientOVPN(c.ConfigDir, name)
	if err != nil {
		if os.IsNotExist(err) {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = JSONResponse{Status: "error", Message: "Certificate not found"}
			c.ServeJSON()
			return
		}
		c.ServeJSONError(err.Error())
		return
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = JSONResponse{Status: "error", Message: "Certificate not found"}
			c.ServeJSON()
			return
		}
		c.ServeJSONError(err.Error())
		return
	}

	filename := name + ".ovpn"
	c.Ctx.Output.Header("Content-Type", "application/x-openvpn-profile")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Ctx.Output.Body(data)
}

// Create creates a new client certificate.
// @Title Create certificate
// @Description Create a new OpenVPN client certificate
// @Param    body     body     controllers.APICreateCertRequest     true      "Certificate parameters"
// @Success 200 request success
// @Failure 400 request failure
// @Failure 401 unauthorized
// @router / [post]
func (c *APICertificatesController) Create() {
	req := APICreateCertRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	if err := c.fillCreateDefaults(&req); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	if err := validateAPICreateCert(req); err != nil {
		c.ServeJSONError(err.Error())
		return
	}

	logs.Info("API: Creating certificate: name=%s", req.Name)
	if err := lib.CreateCertificate(
		req.Name, req.StaticIP, req.Passphrase, req.ExpireDays,
		req.Email, req.Country, req.Province,
		strconv.Quote(req.City), strconv.Quote(req.Org), strconv.Quote(req.OrgUnit),
		req.TFAName, req.TFAIssuer,
	); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	c.ServeJSONMessage("Certificate \"" + req.Name + "\" has been created")
}

// Renew renews an existing client certificate.
// @Title Renew certificate
// @Description Renew an existing OpenVPN client certificate
// @Param    name     path     string     true      "Certificate name"
// @Success 200 request success
// @Failure 400 request failure
// @Failure 401 unauthorized
// @router /:name/renew [post]
func (c *APICertificatesController) Renew() {
	name := c.GetString(":name")
	cert, err := lib.FindActiveCertByName(name)
	if err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	localip := ""
	if cert.Details != nil {
		localip = cert.Details.LocalIP
	}

	logs.Info("API: Renewing certificate: name=%s serial=%s", name, cert.Serial)
	if err := lib.RenewCertificate(name, localip, cert.Serial, certTFAName(cert)); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	c.ServeJSONMessage("Certificate \"" + name + "\" has been renewed")
}

// Revoke revokes a client certificate.
// @Title Revoke certificate
// @Description Revoke an OpenVPN client certificate
// @Param    name     path     string     true      "Certificate name"
// @Success 200 request success
// @Failure 400 request failure
// @Failure 401 unauthorized
// @router /:name/revoke [post]
func (c *APICertificatesController) Revoke() {
	name := c.GetString(":name")
	revoked := false
	for i := 0; i < 5; i++ {
		cert, err := lib.FindActiveCertByName(name)
		if err != nil {
			break
		}
		logs.Info("API: Revoking certificate: name=%s serial=%s", name, cert.Serial)
		if err := lib.RevokeCertificate(name, cert.Serial, certTFAName(cert)); err != nil {
			c.ServeJSONError(err.Error())
			return
		}
		revoked = true
	}
	if !revoked {
		c.ServeJSONError(fmt.Sprintf("active certificate for %q not found", name))
		return
	}
	c.ServeJSONMessage("Certificate \"" + name + "\" has been revoked")
}

// Delete permanently removes a revoked client certificate.
// @Title Delete certificate
// @Description Permanently remove a revoked OpenVPN client certificate
// @Param    name     path     string     true      "Certificate name"
// @Success 200 request success
// @Failure 400 request failure
// @Failure 401 unauthorized
// @router /:name [delete]
func (c *APICertificatesController) Delete() {
	name := c.GetString(":name")
	cert, err := lib.FindRevokedCertByName(name)
	if err != nil {
		c.ServeJSONError(err.Error())
		return
	}

	logs.Info("API: Deleting certificate: name=%s serial=%s", name, cert.Serial)
	if err := lib.BurnCertificate(name, cert.Serial, certTFAName(cert)); err != nil {
		c.ServeJSONError(err.Error())
		return
	}
	c.ServeJSONMessage("Certificate \"" + name + "\" has been deleted")
}

func certTFAName(cert *lib.Cert) string {
	if cert.Details == nil {
		return ""
	}
	if cert.Details.TFAName == "" || cert.Details.TFAName == "none" {
		return ""
	}
	return cert.Details.TFAName
}

func (c *APICertificatesController) fillCreateDefaults(req *APICreateCertRequest) error {
	easyRSA := models.EasyRSAConfig{Profile: "default"}
	if err := easyRSA.Read("Profile"); err != nil {
		return err
	}
	clientCfg := models.OVClientConfig{Profile: "default"}
	if err := clientCfg.Read("Profile"); err != nil {
		return err
	}

	if req.ExpireDays == "" {
		req.ExpireDays = strconv.Itoa(easyRSA.EasyRSACertExpire)
	}
	if req.Email == "" {
		req.Email = easyRSA.EasyRSAReqEmail
	}
	if req.Country == "" {
		req.Country = easyRSA.EasyRSAReqCountry
	}
	if req.Province == "" {
		req.Province = easyRSA.EasyRSAReqProvince
	}
	if req.City == "" {
		req.City = easyRSA.EasyRSAReqCity
	}
	if req.Org == "" {
		req.Org = easyRSA.EasyRSAReqOrg
	}
	if req.OrgUnit == "" {
		req.OrgUnit = easyRSA.EasyRSAReqOu
	}
	if req.TFAIssuer == "" {
		req.TFAIssuer = clientCfg.TFAIssuer
	}
	return nil
}

func validateAPICreateCert(req APICreateCertRequest) error {
	valid := validation.Validation{}
	if ok, err := valid.Valid(&req); err != nil {
		return err
	} else if !ok {
		return validationErr(valid)
	}
	return nil
}

func validationErr(valid validation.Validation) error {
	for _, e := range valid.Errors {
		return fmt.Errorf("%s: %s", e.Key, e.Message)
	}
	return fmt.Errorf("validation failed")
}

func apiSaveClientOVPN(configDir, name string) (string, error) {
	cfg := clientconfig.New()
	keysPath := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/issued")
	keysPathCa := filepath.Join(state.GlobalCfg.OVConfigPath, "pki")

	ovClientConfig := &models.OVClientConfig{Profile: "default"}
	if err := ovClientConfig.Read("Profile"); err != nil {
		return "", err
	}
	cfg.ServerAddress = ovClientConfig.ServerAddress
	cfg.OpenVpnServerPort = ovClientConfig.OpenVpnServerPort
	cfg.AuthUserPass = ovClientConfig.AuthUserPass
	cfg.ResolveRetry = ovClientConfig.ResolveRetry
	cfg.OVClientUser = ovClientConfig.OVClientUser
	cfg.OVClientGroup = ovClientConfig.OVClientGroup
	cfg.PersistTun = ovClientConfig.PersistTun
	cfg.PersistKey = ovClientConfig.PersistKey
	cfg.RemoteCertTLS = ovClientConfig.RemoteCertTLS
	cfg.RedirectGateway = ovClientConfig.RedirectGateway
	cfg.Proto = ovClientConfig.Proto
	cfg.Auth = ovClientConfig.Auth
	cfg.Cipher = ovClientConfig.Cipher
	cfg.Device = ovClientConfig.Device
	cfg.AuthNoCache = ovClientConfig.AuthNoCache
	cfg.TlsClient = ovClientConfig.TlsClient
	cfg.Verbose = ovClientConfig.Verbose
	cfg.CustomConfOne = ovClientConfig.CustomConfOne
	cfg.CustomConfTwo = ovClientConfig.CustomConfTwo
	cfg.CustomConfThree = ovClientConfig.CustomConfThree

	ca, err := os.ReadFile(filepath.Join(keysPathCa, "ca.crt"))
	if err != nil {
		return "", err
	}
	cfg.Ca = string(ca)

	ta, err := os.ReadFile(filepath.Join(keysPathCa, "ta.key"))
	if err != nil {
		return "", err
	}
	cfg.Ta = string(ta)

	cert, err := os.ReadFile(filepath.Join(keysPath, name+".crt"))
	if err != nil {
		return "", err
	}
	cfg.Cert = string(cert)

	keysPathKey := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/private")
	key, err := os.ReadFile(filepath.Join(keysPathKey, name+".key"))
	if err != nil {
		return "", err
	}
	cfg.Key = string(key)

	serverConfig := models.OVConfig{Profile: "default"}
	_ = serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "clients", name+".ovpn")
	if err := SaveToFile(filepath.Join(configDir, "openvpn-client-config.tpl"), cfg, destPath); err != nil {
		logs.Error(err)
		return "", err
	}
	return destPath, nil
}
