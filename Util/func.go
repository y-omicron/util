package Util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/projectdiscovery/httpx/runner"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
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

func HttpXFileVerify(isFile bool, TagName string, Proxy string, Thread int, out chan<- []string) (ret int) {
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
		ret += 1
		out <- regExp.FindStringSubmatch(fc.Text())[1:]
	}
	return ret
}
func UnixShell(s string) (string, error) {
	//这里是一个小技巧, 以 '/bin/bash -c xx' 的方式调用shell命令, 则可以在命令中使用管道符,组合多个命令
	cmd := exec.Command("/bin/sh", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out //把执行命令的标准输出定向到out
	cmd.Stderr = &out //把命令的错误输出定向到out

	//启动一个子进程执行命令,阻塞到子进程结束退出
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), err
}

// FormatSize 字节的单位转换 保留两位小数
func FormatSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%7.2f  B", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%7.2f KB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%7.2f MB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%7.2f GB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%7.2f TB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("Size >> TB")
	}
}
