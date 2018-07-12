package windows

import (
	"log"
	"syscall"
	"unsafe"
	//"fmt"
	"fmt"
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
	CreaEnumerateA   = dlladvapi32.NewProc("CredEnumerateA")
)

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func CredManModuleStart() {
	var count uint64
	var creds **Creds
	CreaEnumerateA.Call(0, 0, uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&creds)))


	test := (*[512]Creds)(unsafe.Pointer(*creds))

	//log.Println(*creds)
	fmt.Println(count)
	/*ar test []Creds
	test = *creds*/

//	log.Println((*creds)[0])


	var time = (uint64(test[0].LastWritten.HighDateTime)<<32 )+uint64(test[0].LastWritten.LowDateTime)

	log.Printf("time : %d",time)
	log.Println(test[0])
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
			fmt.Print(string([]byte{p[j]}))
		}
		fmt.Println()
	}
	//log.Println(test[0])

	//log.Println(test[1])
	//fmt.Println(test[0])
	//fmt.Println(test[1])
	log.Println(GetLastError())

}

func CredmanExtractDataRun() {
	log.Println("Credman start")
	CredManModuleStart()
}
