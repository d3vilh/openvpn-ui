package controllers

import (
	"bufio"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/d3vilh/openvpn-ui/models"
)

type LogsController struct {
	BaseController
}

func (c *LogsController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
}

func (c *LogsController) Get() {
	c.TplName = "logs.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Logs",
	}

	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")

	if err := settings.Read("OVConfigPath"); err != nil {
		logs.Error(err)
		return
	}

	fName := settings.OVConfigPath + "/log/openvpn.log"
	file, err := os.Open(fName)
	if err != nil {
		logs.Error(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logs []string
	for scanner.Scan() {
		line := scanner.Text()
		//	if strings.Index(line, " MANAGEMENT: ") == -1 {
		if !strings.Contains(line, " MANAGEMENT: ") {
			logs = append(logs, strings.Trim(line, "\t"))
		}
	}
	start := len(logs) - 300 // :P
	if start < 0 {
		start = 0
	}
	c.Data["logs"] = logs[start:]
	//c.Data["logs"] = reverse(logs[start:])
}

//func reverse(lines []string) []string {
//	for i := 0; i < len(lines)/2; i++ {
//		j := len(lines) - i - 1
//		lines[i], lines[j] = lines[j], lines[i]
//	}
//	return lines
//}
