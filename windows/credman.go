package windows

import (
	"log"
	"syscall"
	"unsafe"
)

type lifetime struct {
	LowDateTime uint32
	HighDateTime uint32
}

type CredAttr struct {
	Keyword *[]byte
	Flags uint32
	Valuesize uint32
	Value *[]byte
}

type Creds struct {
	FLAG uint32
	TYPE uint32
	TargetName *[]byte
	Comment *[]byte
	LastWritten lifetime
	CredentialBlobSize uint32
	CredentialBlob *[]byte
	Persist uint32
	AttrCount uint32
	Attributes CredAttr
	TargetAlias *[]byte
	Username *[]byte
}

var (

	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetLastError       = modkernel32.NewProc("GetLastError")
	dlladvapi32  = syscall.NewLazyDLL("Advapi32.dll")
	CreaEnumerateA   = dlladvapi32.NewProc("CredEnumerateA")
)

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func CredManModuleStart(){
	var count uint64
	var creds **Creds
	CreaEnumerateA.Call(0, 0, uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&creds)))
	log.Println(&creds)
	log.Println(GetLastError())

}

func CredmanExtractDataRun(){
	log.Println("Credman start")
	CredManModuleStart()
}
