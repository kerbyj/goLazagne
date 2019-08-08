package sysadmin

import (
	"fmt"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/types"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"strings"
)

//retrieves username, host, port
func retrieveRest(sessionInfo string) (string, string, string) {
	piece := sessionInfo[strings.Index(sessionInfo, "%")+1:]
	//extracts port
	pattern1 := `(^([A-Za-z0-9\.\-_]*%)?)|((%[A-Za-z0-9:,\\<>\?\.\-_# ]*)*)`
	re := regexp.MustCompile(pattern1)
	port := re.ReplaceAllString(piece, "")
	//extracts username
	pattern2 := `(^([A-Za-z0-9\.\-_]*%){2})|((%[A-Za-z0-9:,\\<>\?\.\-_# ]*)*)`
	reg := regexp.MustCompile(pattern2)
	userName := reg.ReplaceAllString(piece, "")
	host := piece[:strings.Index(piece, "%")]
	return userName, host, port
}

//returns key location & key
func retrieveKey(sessionInfo string) (string, string) {
	pattern := `(^([A-Za-z0-9_ \?\.\-#]*%)+)|(%[A-Za-z0-9,<> \.\-#]*)`
	re := regexp.MustCompile(pattern)
	keyLocation := re.ReplaceAllString(sessionInfo, "")
	userPath, err := os.UserHomeDir()
	if err != nil {
		return "", ""
	}
	//get system drive from env variable, otherwise make it to 'C'
	systemDrive := os.Getenv("SystemDrive")
	if len(systemDrive) < 0 {
		systemDrive = "C"
	}
	//_ProfileDir_ and _CurrentDrive_ are ini's '%USERPROFILE%' and 'C'
	keyLocation = strings.Replace(keyLocation, "_ProfileDir_", userPath, 1)
	keyLocation = strings.Replace(keyLocation, "_CurrentDrive_:", systemDrive, 1)
	_, err = os.Open(keyLocation)
	if err != nil {
		return "", ""
	}
	key := common.ReadKey(keyLocation)
	if key == nil {
		return "", ""
	}
	if common.OpensshKeyCheck(key) || common.PpkKeyCheck(key) {
		return keyLocation, string(key)
	} else {
		return "", ""
	}
}

//retrieves info from ini config file and returns array of hosts, users, key locations, ports and keys
func MobaExtractDataRun() ([]types.MobaData, error) {
	userPath, err := os.UserHomeDir()
	f, err := os.Open(userPath + "/Documents/MobaXterm/MobaXterm.ini")

	if err != nil {
		return nil, err
	}
	defer f.Close()
	fileStat, _ := f.Stat()
	size := int(fileStat.Size())
	data := make([]byte, size)
	_, err = f.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Error reading moba config file")
	}
	iniRaw, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
		Insensitive:              true,
		AllowShadows:             false,
	}, userPath+"/Documents/MobaXterm/MobaXterm.ini")
	if err != nil {
		fmt.Println("Error with ini reader: ", err)
		return nil, err
	}
	sections := iniRaw.SectionStrings()
	var sessionsInfo []types.MobaData
	var tempSessInfo types.MobaData
	//sessionsInfo := [6]types.MobaData{}
	var info = make(map[string][]string)
	for _, value := range sections {
		if strings.Contains(value, "bookmarks") {
			section := iniRaw.Section(value)
			//iterate through bookmark's field values
			//skip folders without sessions( > 2 check)
			for i := 2; i < len(section.KeyStrings()) && len(section.KeyStrings()) > 2; i++ {
				session := section.Key(section.KeyStrings()[i])
				keyLocation, key := retrieveKey(session.Value())
				if keyLocation == "" {
					continue
				}
				if userName, host, port := retrieveRest(session.Value()); host != "" {
					tempSessInfo.HostName = host
					tempSessInfo.User = userName
					tempSessInfo.KeyLocation = keyLocation
					tempSessInfo.Port = port
					tempSessInfo.Key = key
				}
				sessionsInfo = append(sessionsInfo, tempSessInfo)
			}
		}
	}
	if len(info) < 0 {
		return sessionsInfo, nil
	}
	return sessionsInfo, nil
}
