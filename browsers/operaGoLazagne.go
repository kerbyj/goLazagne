package browsers


import (
	"os"
	"fmt"
	"database/sql"
	"goLaZagne/common"
)

var (
	operaPathsUserData = []string{
		common.AppData +"\\Opera Software\\Opera Stable",
	}
	temporaryDbName = common.RandStringRunes(10)
)

func OperaModuleStart(path string) ([]string, bool){
	if _, err := os.Stat(path); err == nil {
		dbPath := fmt.Sprintf("%s\\Login data", path)

		if _, err := os.Stat(dbPath); err == nil {
			err := common.CopyFile(dbPath, temporaryDbName)
			if err != nil{
				//log.Panic(err.Error())
				return nil, false
			}
		}

		db, err := sql.Open("sqlite3", temporaryDbName)
		if err != nil {
			return nil, false
		}
		rows, err := db.Query("SELECT action_url, username_value, password_value FROM logins")
		var actionUrl, username, password string
		var data []string
		for rows.Next(){
			rows.Scan(&actionUrl, &username, &password)
			//fmt.Printf("%s %s - %s", actionUrl, username, common.Win32CryptUnprotectData(password, false))
			data = append(data, fmt.Sprintf("%s %s %s", actionUrl, username, common.Win32CryptUnprotectData(password, false)))
		}
		return data, true
	}
	return nil, false
}

func OperaExtractDataRun() common.ExtractDataResult{
	var Result common.ExtractDataResult
	var EmptyResult = common.ExtractDataResult{false,Result.Data}
	for i:=range operaPathsUserData {
		if _, err := os.Stat(operaPathsUserData[i]); err == nil {
			var data, success = OperaModuleStart(operaPathsUserData[i])
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