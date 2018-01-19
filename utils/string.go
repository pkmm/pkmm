package utils

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"os"
	"github.com/astaxie/beego"
	"strings"
	"encoding/hex"
	"io"
	"encoding/base64"
	"crypto/rand"
	"runtime"
	"path"
)

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		beego.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

/// 路径
func GetSourceCodePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func GetExecPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}