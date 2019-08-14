# goLazagne

---
## ⚠️ Disclaimer

1. All information provided in this project is for educational purposes only and can not be used for law violation or personal gain.
2. The authors of this project are not responsible for any possible harm caused by the materials of this project.
3. All information in this repository is intended for development of audit tools and help preventing the hack attacks.
4. We believe only in ethical hacking.
---

## Description

The **goLazagne** is an open source library for golang used to retrieve passwords stored on local computer.

Inspired by AlessandroZ [LaZagne](https://github.com/AlessandroZ/LaZagne) project.

## Supported features

* Browsers
	* Chromium-based
	* Mozilla Firefox
	* Internet Explorer and Edge

* Mail
    * Thunderbird

* Windows
    * Credential Manager
    
* SysAdmin tools (pre alpha, need more tests)
    * Mobaxterm - user, host:port and associated key
    * Putty - user, host:port and associated key
    * Filezilla - user, host:port and password (encrypted if master password isset)
    * Openssh

* WiFi passwords
	
## ToDo (sorted by priority level)

- [ ] WPA2 Enterprise. The main difficulty is that we need an privilege escalation. Read more in zc00l [research](https://0x00-0x00.github.io/research/2018/11/06/Recovering-Plaintext-Domain-Credentials-From-WPA2-Enterprise-on-a-compromised-host.html).
- [ ] Windows vault
- [ ] Outlook
- [ ] Git 
    
## Special thanks

* [Nikolay Edigaryev](https://github.com/edigaryev) for [credman](https://github.com/kerbyj/goLazagne/blob/master/windows/credman.go) module

## Example

```go
package main

import (
    "github.com/kerbyj/goLazagne"
)

func main() {

    var credentials, _ = goLazagne.ExtractAllData()
    
    println("Browser creds:", len(credentials.BrowserData))
    println("Credman creds:", len(credentials.CredmanData))
    println("Wifi creds:", len(credentials.WifiData))
    
    println("\nEnumerating filesystem. Please wait")
    
    var interestingFiles = []string{
    			"ovpn",
    			"ssh",
    		}
    
    var files = goLazagne.ExtractInterestingFiles(interestingFiles)
    for fileN := range files {
        println(files[fileN])
    }

}
```
