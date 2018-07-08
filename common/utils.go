package common

import (
	"io/ioutil"
	"unsafe"
	"syscall"
	"math/rand"
	"os"
)

var (

	UserHome = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	AppData  = os.Getenv("APPDATA")

)

func RandStringRunes(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CopyFile(src string, dst string) error{
	data, err := ioutil.ReadFile(src)
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil{
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

func Win32CryptUnprotectData(cipherText string, entropy bool) string{
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
