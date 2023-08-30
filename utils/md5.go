package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(MD5Encode(data))
}

// 加密
func MakePassword(plainPwd, salt string) string {
	return Md5Encode(plainPwd + salt)
}

func ValidPassword(plainPwd, salt, password string) bool {
	fmt.Println(Md5Encode(plainPwd+salt), password)
	return Md5Encode(plainPwd+salt) == password
}
