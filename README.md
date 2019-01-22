### goLazagne

## Supported feautures

* Browsers
	* Google Chrome 
	* Mozilla Firefox 
	* Opera
	* Yandex Browser

* Mail
    * Thunderbird

* Windows
    * Credential Manager // windows 10 only

* WiFi passwords
	* WPA2-PSK

#### Example

```go
package main

import (
    "github.com/kerbyj/goLazagne"
)

func main() {

    var credentials, _ = GoLazagne.ExtractAllData()
    
    println("Browser creds:", len(credentials.BrowserData))
    println("Credman creds:", len(credentials.CredmanData))
    println("Wifi creds:", len(credentials.WifiData))
    
    println("\nEnumerating filesystem. Please wait")
    var files = GoLazagne.ExtractInterstingFiles()
    for fileN := range files {
        println(files[fileN])
    }

}
```
