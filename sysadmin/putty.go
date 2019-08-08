package sysadmin

//TODO
//extract ports as well
//HKCU/Software/*Sessions/#SESSNAME#/PortNumber (it's in hex)
import (
	"fmt"
	"github.com/kerbyj/goLazagne/common"
	"github.com/kerbyj/goLazagne/types"
	"golang.org/x/sys/windows/registry"
	"log"
	"os/exec"
	"strings"
)

func hostNameExtractor(k registry.Key) string {
	hostName, _, err := k.GetStringValue("HostName")
	if err != nil {
		log.Println("Error extracting hostname: ", err)
	}
	return hostName
}

func userNameExtractor(k registry.Key) string {
	userName, _, err := k.GetStringValue("UserName")
	//we can work w/o username
	if err != nil {
		log.Println("Error extracting username: ", err)
	}
	return userName
}

func keyExtractor(k registry.Key) string {
	privKeyPath, _, err := k.GetStringValue("PublicKeyFile")
	if err != nil {
		log.Println("Error extracting private key location: ", err)
		return ""
	}
	key := common.ReadKey(privKeyPath)
	if key != nil && (common.PpkKeyCheck(key) || common.OpensshKeyCheck(key)) {
		return string(key)
	} else {
		return ""
	}
}

//extracts user, key, hostname
func puttyInfo(pathToSession string) (string, string, string) {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		pathToSession, registry.QUERY_VALUE)
	if err != nil {
		log.Println("Error opening registry: ", err)
		return "", "", ""
	}
	hostName := hostNameExtractor(k)
	userName := userNameExtractor(k)
	key := keyExtractor(k)
	return hostName, userName, key
}

//extract Putty's username, hostname & key location from registry
func PuttyExtractor() ([]types.PuttyData, error) {
	var keys []types.PuttyData
	//get the sessions hives' names
	output, err := exec.Command("powershell",
		"reg query HKCU\\Software\\SimonTatham\\Putty\\Sessions").Output()
	if err != nil {
		return keys, fmt.Errorf("powershell failed, putty may not be installed")
	}
	if len(output) < 0 {
		fmt.Print("len(output) < 0 ..")
		return keys, err
	}
	out := strings.Split(string(output), "\r\n")
	out = out[1 : len(out)-1]
	for i := range out {
		out[i] = out[i][18:]
		hostName, userName, key := puttyInfo(out[i])
		if key != "" {
			temp := types.PuttyData{HostName: hostName, UserName: userName, Key: key}
			keys = append(keys, temp)
		}
	}
	return keys, err

}

func PuttyExtractDataRun() ([]types.PuttyData, error) {
	info, err := PuttyExtractor()
	if err != nil {
		return info, err
	}
	return info, nil
}
