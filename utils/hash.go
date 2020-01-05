package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

func HASH(msg, salt string) string {
	m := md5.New()
	s := sha256.New()
	s.Write([]byte(msg + salt))
	m.Write([]byte(hex.EncodeToString(s.Sum(nil)) + salt))
	return hex.EncodeToString(m.Sum(nil))
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
