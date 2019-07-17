package wifi

import (
	"github.com/aglyzov/charmap"
	"github.com/kerbyj/goLazagne/common"
	"os/exec"
	"strings"
	"syscall"
)

func execCommand(command string, params []string) string {
	cmd_li := exec.Command(command, params...)
	cmd_li.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //run CMD in hidden mode
	output, _ := cmd_li.Output()
	if output != nil && len(output) > 0 {
		output = charmap.CP866_to_UTF8(output)
	}
	return string(output)
}

/*
	Function for WiFi credentials extracting. No support for WPA2 Enterprise (see README).
 */
func WifiExtractDataRun() common.ExtractCredentialsNamePass {

	/**
		For WiFI credentials extract we use system utility "netsh"

		1. Get profile names (SSID) - `netsh wlan show profiles`
		2. Get saved password for wifi access point - `netsh wlan show profile SSID key=clear`
			key=clear parameter used to display the password in clear text
	 */

	params := []string{
		"wlan",
		"show",
		"profiles",
	}

	var output = execCommand("netsh", params)
	var lines = strings.Split(output, "\r\n")
	var users []string

	for i := range lines {
		if strings.Contains(lines[i], "All User") { //TODO check in multiple languages
			users = append(users, strings.TrimSpace(strings.Split(lines[i], ":")[1]))
		}
	}

	var Result common.ExtractCredentialsNamePass
	var data []common.NamePass
	for i := 0; i < len(users); i++ {
		var paramWifi = []string{
			"wlan",
			"show",
			"profile",
			users[i],
			"key=clear",
		}

		var output = execCommand("netsh", paramWifi)
		var lines = strings.Split(output, "\r\n")

		for j := range lines {
			if strings.Contains(lines[j], "Key Content") { //TODO check in multiple languages
				var (
					dataAdd = common.NamePass{
						users[i],
						strings.TrimSpace(strings.Split(lines[j], ":")[1]),
					}
				)
				data = append(data, dataAdd)
			}
		}
		if len(data) == 0 {
			Result.Success = false
			return Result
		}
	}
	Result.Data = data
	Result.Success = true
	return Result
}
