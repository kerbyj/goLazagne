package goLazagne

import (
	"github.com/kerbyj/goLazagne/browsers"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/filesystem"
	"github.com/kerbyj/goLazagne/sysadmin"
	"github.com/kerbyj/goLazagne/types"
	"github.com/kerbyj/goLazagne/wifi"
	"github.com/kerbyj/goLazagne/windows"
)

/**
Common function for work with browsers. Just call and function return all saved passwords in chromium browsers and firefox
*/
func ExtractBrowserCredentials() ([]common.UrlNamePass, int) {
	var AllBrowsersData []common.UrlNamePass
	if resultChrome := browsers.ChromeExtractDataRun(); resultChrome.Success {
		AllBrowsersData = append(AllBrowsersData, resultChrome.Data...)
	}

	if resultMozilla := browsers.MozillaExtractDataRun("browser"); resultMozilla.Success {
		AllBrowsersData = append(AllBrowsersData, resultMozilla.Data...)
	}

	if resultInternetExplorer := browsers.InternetExplorerExtractDataRun(); resultInternetExplorer.Success {
		AllBrowsersData = append(AllBrowsersData, resultInternetExplorer.Data...)
	}

	return AllBrowsersData, len(AllBrowsersData)
}

/*
	Function that check saved credentials in chromium based browsers
*/
func ExtractChromiumCredentials() common.ExtractCredentialsResult {
	return browsers.ChromeExtractDataRun()
}

/**
Function that check saved credentials in firefox browser and thunderbird

targetType parameter: "mail" parameter value is used for Thunderbird processing, all others for Firefox-based software
*/
func ExtractFirefoxCredentials() common.ExtractCredentialsResult {
	return browsers.MozillaExtractDataRun("browser")
}

// Function that check saved credentials in internet explorer and edge
func ExtractIECredentials() common.ExtractCredentialsResult {
	return browsers.InternetExplorerExtractDataRun()
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

type SysadminData struct {
	MobaXTerm []types.MobaData      `json:"moba_x_term"`
	OpenSsh   types.OpensshData     `json:"open_ssh"`
	Putty     []types.PuttyData     `json:"putty"`
	Filezilla []types.FileZillaData `json:"filezilla"`
	Count     int                   `json:"count"`
}

func ExtractSysadminData() SysadminData {
	var outdata SysadminData

	mobaData, errExtractMobaData := sysadmin.MobaExtractDataRun()
	openssh, errOpenSsh := sysadmin.OpensshExtractDataRun()
	putty, errPutty := sysadmin.PuttyExtractDataRun()
	filezilla, errFileZilla := sysadmin.FilezillaExtractDataRun()
	outdata.Count = 0

	if errExtractMobaData == nil {
		outdata.MobaXTerm = mobaData
		outdata.Count += len(mobaData)
	}

	if errOpenSsh == nil {
		outdata.OpenSsh = openssh
		outdata.Count += len(openssh.Hosts)
	}

	if errPutty == nil {
		outdata.Putty = putty
		outdata.Count += len(putty)
	}

	if errFileZilla == nil {
		outdata.Filezilla = filezilla
		outdata.Count += len(filezilla)
	}

	return outdata
}

type AllDataStruct struct {
	WifiData     []common.NamePass    `json:"wifi"`
	BrowserData  []common.UrlNamePass `json:"browser"`
	CredmanData  []common.UrlNamePass `json:"credman"`
	SysadminData SysadminData         `json:"sysadmin_data"`
}

//Function in "give me all" style. The function will return everything that the program can extract from OS.
func ExtractAllData() (AllDataStruct, int) {
	var wifiData, lengthWiFiData = ExtractWifiData()
	var browserData, lengthBrowserData = ExtractBrowserCredentials()
	var credmanData, lengthCredmanData = ExtractCredmanData()
	var sysadminData = ExtractSysadminData()

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
	outDataStruct.SysadminData = sysadminData

	return outDataStruct, lengthCredmanData + lengthBrowserData + lengthWiFiData
}
