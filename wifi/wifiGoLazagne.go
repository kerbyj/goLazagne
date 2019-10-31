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

	for _, line := range lines {
		if strings.Contains(line, "Все профили") || strings.Contains(line, "All profile") {
			users = append(users, strings.TrimSpace(strings.Split(line, ":")[1]))
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

		for i, line := range lines {
			if strings.Contains(line, "Содержимое ключа") || strings.Contains(line, "Key content")  { //todo verify this part of code
				var (
					dataAdd = common.NamePass{
						Name: users[i],
						Pass: strings.TrimSpace(strings.Split(line, ":")[1]),
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
