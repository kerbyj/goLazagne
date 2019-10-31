package sysadmin

//TODO
//known_hosts(?)
//move common to common
import (
	"fmt"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/types"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func retrieveHostname() []string {
	f, err := os.Open(common.UserHome + "/.SSH/known_hosts")
	if err != nil {
		return nil
	}
	stat, err := f.Stat()
	if err != nil {
		return nil
	}
	size := int(stat.Size())
	contents := make([]byte, size)
	_, err = f.Read(contents)
	//error reading file
	if err != nil {
		log.Println("Error reading file: ", err)
		return nil
	}
	pattern := `(^([A-Za-z0-9\.\:\[\]\,\@]*? ){1})`
	re := regexp.MustCompile(pattern)
	match := re.FindAllString(string(contents), -1)
	return match
}

//Retrieve ssh private keys
func extractUnprotectedPrivKeys() []string {
	var keys []string
	//default key files path %USERPROFILE%/.SSH
	keyFilesLocation := common.UserHome + "/.SSH"
	folder, err := os.Stat(keyFilesLocation)

	//check if key files path is present
	if err == nil && folder.IsDir() {
		err = filepath.Walk(keyFilesLocation, func(path string, info os.FileInfo, err error) error {
			//failed accessing the file
			if err != nil {
				return err
			}
			//skip check for directories
			if info.IsDir() {
				return nil
			}
			key := common.ReadKey(path)
			if key == nil {
				return nil
			}
			if common.OpensshKeyCheck(key) || common.PpkKeyCheck(key) {
				keys = append(keys, string(key))
			}
			return nil
		})
	}
	return keys
}

func OpensshExtractDataRun() (types.OpensshData, error) {
	var data types.OpensshData
	data.Hosts = retrieveHostname()

	data.Keys = extractUnprotectedPrivKeys()

	if len(data.Keys) == 0 {
		return data, fmt.Errorf("Nothing found.")
	} else {
		return data, nil
	}
}
