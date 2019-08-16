package wifi

import (
	"github.com/kerbyj/goLazagne/common"
	"strings"
)

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
		"netsh",
		"wlan",
		"show",
		"profiles",
	}

	var output = common.ExecCommand("cmd", params)
	var lines = strings.Split(output, "\r\n")
	var users []string

	for i := range lines {
		if strings.Contains(lines[i], "Все профили") { //TODO check in multiple languages
			users = append(users, strings.TrimSpace(strings.Split(lines[i], ":")[1]))
		}
	}

	var Result common.ExtractCredentialsNamePass
	var data []common.NamePass
	for i := 0; i < len(users); i++ {
		var paramWifi = []string{
			"netsh",
			"wlan",
			"show",
			"profile",
			users[i],
			"key=clear",
		}

		var output = common.ExecCommand("cmd", paramWifi)
		var lines = strings.Split(output, "\r\n")

		for j := range lines {
			if strings.Contains(lines[j], "Содержимое ключа") { //TODO check in multiple languages
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
