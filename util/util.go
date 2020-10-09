package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * 将文件大小转换成人类认识的展示方式
 */
func FormatSize(size int64) string {
	fSize, _ := strconv.ParseFloat(strconv.FormatInt(size, 10), 32/64)
	if fSize < 1024 {
		return fmt.Sprintf("%.2f", fSize) + "B"
	}
	fSize = fSize / 1024
	if fSize < 2048 {
		return fmt.Sprintf("%.2f", fSize) + "KB"
	}
	fSize = fSize / 1024
	if fSize < 1024 {
		return fmt.Sprintf("%.2f", fSize) + "MB"
	}
	return fmt.Sprintf("%.2f", fSize/1024) + "GB"
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GenerateToken() string {
	now := time.Now().Unix()
	h := md5.New()
	if _, err := io.WriteString(h, strconv.FormatInt(now, 10)); err != nil {
		return ""
	}
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

func GetClientIp() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addr {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("can not find the client ip address")
}

/**
 * 查询文件（文件件）是否存在
 */
func ExistsFile(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	return err == nil || os.IsExist(err)
}
func ExistsPath(path string) bool {
	f, err := os.Stat(path) //os.Stat获取文件信息
	return (err == nil || os.IsExist(err)) && f.IsDir()
}

/**
 * 返回上传文件的MD5校验码
 */
func Md5File(file multipart.File) string {
	m := md5.New()
	_, err := io.Copy(m, file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	md5s := hex.EncodeToString(m.Sum(nil))
	return md5s
}

/**
 * 异常处理
 */
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
