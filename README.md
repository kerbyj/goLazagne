# goLazagne

---

#### ⚠️ Disclaimer

1. All the information provided in this project is for educational purposes only and can not be used for law violation or personal gain.
2. The authors of this project are not responsible for any possible harm caused by the materials of this project.
3. All information in this repository is intended for development of audit tools and help preventing the hack attacks.
4. We believe only in ethical hacking.
---

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
