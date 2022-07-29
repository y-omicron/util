package Util

import (
	"archive/zip"
	"bytes"
	"math/rand"
	"os"
	"time"
)

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
func RandString(strLen int) string {
	strList := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	aBytes := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		c := strList[r.Intn(26)]
		aBytes[i] = byte(c)
	}
	return string(aBytes)
}
func OpenFileToWrite(filename string, date []byte) {
	wf, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	_, err = wf.Write(date)
	if err != nil {
		panic(err)
	}
	wf.Close()
}

// CreateZipFile 制作Zip文件
func CreateZipFile(FileName string, FileDate []byte) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	ioWriter, err := zipWriter.Create(FileName)
	if err != nil {
		return nil, err
	}
	_, err = ioWriter.Write(FileDate)
	if err != nil {
		return nil, err
	}
	zipWriter.Close()
	return buf.Bytes(), nil
}
