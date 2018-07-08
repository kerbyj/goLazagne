package main

import (
	"goLaZagne/browsers"
	"log"
)

func main() {
	//log.Println("Chrome started")
	//browsers.ChromeExtractDataRun()

	//log.Println("Opera started")
	//browsers.OperaExtractDataRun()

	log.Println("Mozilla started")
	browsers.MozillaExtractDataRun()

	//log.Println("WiFi started")
	//wifi.WifiExtractDataRun()
}