package common

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"io"
	"strings"
	"sync"
)

const TokenKey = "ptoken"

type VerifyInfo struct {
	Id      int
	Account string
	Token   string
}

var verifyInfoPool = &sync.Pool{
	New: func() interface{} {
		return &VerifyInfo{}
	},
}

func getVerifyInfo() *VerifyInfo {
	info := verifyInfoPool.Get().(*VerifyInfo)
	info.Id = 0
	info.Account = ""
	info.Token = ""
	return info
}

func DropVerifyInfo(info *VerifyInfo) {
	verifyInfoPool.Put(info)
}

func GetVerifyInfo(tokenStr string) (*VerifyInfo, error) {
	str, err := base64.StdEncoding.DecodeString(tokenStr)
	if err == nil {
		info := getVerifyInfo()
		decode := gob.NewDecoder(bytes.NewBuffer(str))
		err = decode.Decode(info)
		if err == nil {
			return info, nil
		}
	}
	return nil, err
}

func GenerateVerifyToken(id int, account string) ([]byte, error) {
	info := &VerifyInfo{Id: id, Account: account}
	info.Token = strings.TrimRight(base32.StdEncoding.EncodeToString(GenerateRandomKey(32)), "=")
	info.Token = string(GenerateRandomKey(8))

	buf := new(bytes.Buffer)
	buf.Reset()
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(info); err != nil {
		return nil, err
	} else {
		bs := base64.StdEncoding.EncodeToString(buf.Bytes())
		return []byte(bs), err
	}

}

func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
