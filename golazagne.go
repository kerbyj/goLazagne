package goLazagne

import (
	"github.com/kerbyj/goLazagne/browsers"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/filesystem"
	"github.com/kerbyj/goLazagne/wifi"
	"github.com/kerbyj/goLazagne/windows"
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

func ExtractCredmanData() ([]common.UrlNamePass, int) {
	var windowsResult = windows.CredManModuleStart()
	if windowsResult.Success {
		return windowsResult.Data, len(windowsResult.Data)
	}
	return nil, 0
}

func ExtractInterestingFiles() []string {
	var data = filesystem.FindFiles()
	return data
}

type AllDataStruct struct {
	WifiData    []common.NamePass    `json:"wifi"`
	BrowserData []common.UrlNamePass `json:"browser"`
	CredmanData []common.UrlNamePass `json:"credman"`
}

//Structure for fine tuning and work with goLazagne
type LazagneLockPick struct {
	Browsers          bool
	BrowsersSpecify   []string // Work with "chromium" and "firefox"
	CredentialManager bool
	Wifi              bool
	Verbosity         int
	Result            AllDataStruct
}

func (config LazagneLockPick) ExtractData() {

	// Check for zero values in "Browsers" field
	if config.Browsers == true {
		var browserData, lengthBrowserData = ExtractBrowserCredentials()

		if lengthBrowserData > 0 {
			config.Result.BrowserData = append(config.Result.BrowserData, browserData...)
		}

	}

	// Work with specified browsers
	if len(config.BrowsersSpecify) > 0 {

		for i := range config.BrowsersSpecify {

			if config.BrowsersSpecify[i] == "chromium" {

				var chromeExtract = browsers.ChromeExtractDataRun()

				if chromeExtract.Success == true {
					config.Result.BrowserData = append(config.Result.BrowserData, chromeExtract.Data...)
				}

			} else if config.BrowsersSpecify[i] == "firefox" {

				var mozillaExtract = browsers.MozillaExtractDataRun()

				if mozillaExtract.Success == true {
					config.Result.BrowserData = append(config.Result.BrowserData, mozillaExtract.Data...)
				}

			}

		}

	}
	// End work with  specified browsers

	if config.CredentialManager == true {
		var credmanData, lengthCredmanData = ExtractCredmanData()

		if lengthCredmanData > 0 {
			config.Result.CredmanData = credmanData
		}
	}

	if config.Wifi == true {
		var wifiData, lengthWiFiData = ExtractWifiData()

		if lengthWiFiData > 0 {
			config.Result.WifiData = wifiData
		}

	}
}

//Function to extract all credentials from browsers, wifi passwords, and passwords from windows credential manager
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
