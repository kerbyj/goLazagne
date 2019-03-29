### goLazagne

#### Supported features

* Browsers
	* Chromium-based
	* Mozilla Firefox

* Mail
    * Thunderbird

* Windows
    * Credential Manager

* WiFi passwords
	
#### Roadmap

* WPA2 Enterprise
* Outlook
* Windows vault
* Git
    * Git passwords
    
#### Special thanks

* [Nikolay Edigaryev](https://github.com/edigaryev) for [credman](https://github.com/kerbyj/goLazagne/blob/master/windows/credman.go) module

### Example

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
    var files = GoLazagne.ExtractInterestingFiles()
    for fileN := range files {
        println(files[fileN])
    }

}
```
