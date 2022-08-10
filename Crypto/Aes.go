package Crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

type Aes struct {
	Key []byte
	Iv  []byte
}

func (a *Aes) AESGetKey() {
	keyIv := make([]byte, 0x20)
	_, _ = rand.Read(keyIv)
	a.Key = keyIv[:0x10]
	a.Iv = keyIv[0x10:]
}
func (a *Aes) String() string {
	return base64.StdEncoding.EncodeToString(a.Bytes())
}
func (a *Aes) Bytes() []byte {
	return append(a.Key, a.Iv...)
}

// Padding 对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// UnPadding 对密文删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

// AesCbcEncrypt AEC加密（CBC模式）
func (a *Aes) AesCbcEncrypt(plainText []byte) ([]byte, error) {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	//进行填充
	plainText = Padding(plainText, block.BlockSize())
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, a.Iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText, nil
}

// AesCbcDecrypt AEC解密（CBC模式）
func (a *Aes) AesCbcDecrypt(cipherText []byte) ([]byte, error) {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, a.Iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = UnPadding(plainText)
	return plainText, nil
}

func (a *Aes) DecryptToFile(cipherText []byte, FilePath string) (int, error) {
	// 解密数据
	plainText, err := a.AesCbcDecrypt(cipherText)
	if err != nil {
		return 0, err
	}
	// 写入文件
	return len(plainText), os.WriteFile(FilePath, plainText, 0644)
}

func (a *Aes) EncryptToBase64(plainText []byte) (string, error) {
	cbcEncrypt, err := a.AesCbcEncrypt(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cbcEncrypt), nil
}

func TestAes() {
	var SelfAes = Aes{}
	var Plain = []byte("hello world!")
	var Cipher []byte
	var err error
	SelfAes.AESGetKey()
	Cipher, err = SelfAes.AesCbcEncrypt(Plain)
	if err != nil {
		panic(err)
	}
	Plain, err = SelfAes.AesCbcDecrypt(Cipher)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Cipher: %v\nPlain: %s\n", Cipher, Plain)
}
