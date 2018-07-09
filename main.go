package main

import (
	"goLaZagne/wifi"
	"encoding/json"
	"goLaZagne/common"
	"goLaZagne/browsers"
	"log"
	"strings"
)

type SuccessResult struct {
	App string
	Data []common.CredentialsData
}

func packData(result common.ExtractDataResult, name string) []byte{
	var dataForMarshal = SuccessResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func main() {
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

	var BrowsersData = common.ExtractDataResult{false, common.RemoveDuplicates(AllBrowsersData)}
	var data = packData(BrowsersData, "browsers")
	AllData = append(AllData, string(data))

	var resultWifi = wifi.WifiExtractDataRun()
	if resultWifi.Success{
		var data = packData(resultWifi, "wifi")
		AllData = append(AllData, string(data))
	}
	log.Println("["+strings.Join(AllData, ",")+"]")

}