package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
)

//sha1 密码加密
func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum(nil))
}

//MD5 加密
func Md5(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	return hex.EncodeToString(md5.Sum(nil))
}

// Base64加密
func base64Encode(data string) string {
	strbytes := []byte(data)
	encoded := base64.StdEncoding.EncodeToString(strbytes)
	return encoded
}

// Base64解密
func base64Decode(data string) string {
	if decoded, err := base64.StdEncoding.DecodeString(data); err != nil {
		decodestr := string(decoded)
		return decodestr
	}
	return ""
}

// CRC32加密
func CRC32(str string) int {
	return int(crc32.ChecksumIEEE([]byte(str)))
}
