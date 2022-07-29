package Util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"github.com/projectdiscovery/httpx/runner"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
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

func HttpXFileVerify(isFile bool, TagName string, Proxy string, Thread int, out chan<- []string) {
	inputFile := "tmp.ipAddr"
	if !isFile {
		_ = os.WriteFile(inputFile, []byte(TagName), 0644)
		defer os.RemoveAll(inputFile)
	} else {
		inputFile = TagName
	}

	outputFile := "ipAddr.tmp"

	options := runner.Options{
		Methods:         "GET",
		InputFile:       inputFile,
		ExtractTitle:    true,
		StatusCode:      true,
		FollowRedirects: true,
		MaxRedirects:    10,
		RandomAgent:     true,
		Timeout:         10,
		Output:          outputFile,
		Threads:         Thread,
	}
	if Proxy != "" {
		options.HTTPProxy = Proxy
	}
	if err := options.ValidateOptions(); err != nil {
		panic(err)
	}
	HttpXRunner, err := runner.New(&options)
	if err != nil {
		panic(err)
	}
	defer HttpXRunner.Close()

	HttpXRunner.RunEnumeration()
	FileDate, _ := ioutil.ReadFile(outputFile)
	regExp, _ := regexp.Compile(`\x1B\[\d+m`)
	FileDate = regExp.ReplaceAll(FileDate, []byte(""))

	_ = os.WriteFile(outputFile, FileDate, 0644)
	defer os.RemoveAll(outputFile)

	rf, err := os.OpenFile(outputFile, os.O_RDWR, 0644) //以读写方式打开文件
	if nil != err {
		panic(err)
	}
	defer rf.Close()
	fc := bufio.NewScanner(rf) //按行读取文件内容
	regExp, _ = regexp.Compile(`(?U)^(.*) \[(.*)] \[(.*)]`)
	for fc.Scan() {
		out <- regExp.FindStringSubmatch(fc.Text())[1:]
	}
}
