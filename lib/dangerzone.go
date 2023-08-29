package lib

import (
	"fmt"
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
	"github.com/d3vilh/openvpn-ui/state"
)

func DeletePKI() error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && "+
				"./remove-pki.sh"))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func InitPKI() error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && "+
				"./generate_ca_and_server_certs.sh"))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}
