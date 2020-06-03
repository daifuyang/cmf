package util

import (
	"crypto/md5"
	"encoding/hex"

)

var AuthCode *string

func GetMd5(s string) string {
	h := md5.New()
	h.Write([]byte(*AuthCode + s))
	return hex.EncodeToString(h.Sum(nil))
}

