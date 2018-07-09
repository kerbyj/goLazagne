package main

import (
	"goLaZagne/wifi"
	"encoding/json"
	"goLaZagne/common"
	"goLaZagne/browsers"
	"goLaZagne/windows"
)

type SuccessResult struct {
	App string
	Data []string
}

func packData(result common.ExtractDataResult, name string) []byte{
	var dataForMarshal = SuccessResult{name, result.Data}
	var returning, _ = json.Marshal(dataForMarshal)
	return returning
}

func main() {
	var AllBrowsersData []string

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
	var _ = packData(BrowsersData, "browsers")

	//log.Println(string(data))
	var resultWifi = wifi.WifiExtractDataRun()
	if resultWifi.Success{
		var _ = packData(resultWifi, "wifi")
		//log.Println(string(data))
	}

	windows.CredmanExtractDataRun()

}