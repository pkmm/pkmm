package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Today() int{
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Day()
}