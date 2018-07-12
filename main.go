package main

import (
	"goLaZagne/wifi"
	"encoding/json"
	"goLaZagne/common"
	"goLaZagne/browsers"
	"goLaZagne/windows"
	"log"
	"strings"
)

type SuccessCredentialsResult struct {
	App string
	Data []common.CredentialsData
}

type SuccessWifiResult struct {
	App string
	Data []common.WifiData
}

func packBrowsersData(result common.ExtractCredentialsResult, name string) []byte{
	var dataForMarshal = SuccessCredentialsResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func packWifiData(result common.ExtractWifiData, name string) []byte{
	var dataForMarshal = SuccessWifiResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func main() {


	windows.CredmanExtractDataRun()

	return


	var AllBrowsersData []common.CredentialsData
	var AllData []string
	if resultChrome := browsers.ChromeExtractDataRun(); resultChrome.Success{
		AllBrowsersData = append(AllBrowsersData, resultChrome.Data...)
	}
	if resultOpera := browsers.OperaExtractDataRun(); resultOpera.Success{
		AllBrowsersData = append(AllBrowsersData, resultOpera.Data...)
	}
	if resultMozilla := browsers.MozillaExtractDataRun(); resultMozilla.Success {
		AllBrowsersData = append(AllBrowsersData, resultMozilla.Data...)
	}

	var BrowsersData = common.ExtractCredentialsResult{false, common.RemoveDuplicates(AllBrowsersData)}
	var data = packBrowsersData(BrowsersData, "browsers")
	AllData = append(AllData, string(data))

	var resultWifi = wifi.WifiExtractDataRun()
	if resultWifi.Success{
		var data = packWifiData(resultWifi, "wifi")
		AllData = append(AllData, string(data))
	}
	log.Println("["+strings.Join(AllData, ",")+"]")

}