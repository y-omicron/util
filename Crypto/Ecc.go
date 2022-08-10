package Crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

type Ecc struct {
	EcdsaPrivateKey *ecdsa.PrivateKey
	EcdsaPublicKey  *ecdsa.PublicKey
	EciesPrivateKey *ecies.PrivateKey
	EciesPublicKey  *ecies.PublicKey
}

// EccGetKey 生成一个ecc(ecdsa)私钥, 并且将其转换为ecc(ecies)私钥
func (e *Ecc) EccGetKey() error {
	// 初始化椭圆曲线
	PublicCurve := crypto.S256()

	// 随机挑选基点,生成私钥
	p, err := ecdsa.GenerateKey(PublicCurve, rand.Reader)
	if err != nil {
		return err
	}
	e.EcdsaPrivateKey = p

	e.EcdsaPublicKey = &e.EcdsaPrivateKey.PublicKey
	// 将标准包生成私钥转化为ecies私钥
	e.EciesPrivateKey = ecies.ImportECDSA(p)
	e.EciesPublicKey = &e.EciesPrivateKey.PublicKey
	return nil
}

// ECCEncrypt ecc(ecies)加密
func (e *Ecc) ECCEncrypt(pt []byte) ([]byte, error) {
	ct, err := ecies.Encrypt(rand.Reader, e.EciesPublicKey, pt, nil, nil)
	return ct, err
}

// ECCDecrypt ecc(ecies)解密
func (e *Ecc) ECCDecrypt(ct []byte) ([]byte, error) {
	pt, err := e.EciesPrivateKey.Decrypt(ct, nil, nil)
	return pt, err
}

func (e *Ecc) EccSign(pt []byte) (sign []byte, err error) {
	// 根据明文plaintext和私钥，生成两个big.Ing
	r, s, err := ecdsa.Sign(rand.Reader, e.EcdsaPrivateKey, pt)
	if err != nil {
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// 将r，s合并（以“+”分割），作为签名返回
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(`+`))
	b.Write(ss)
	return b.Bytes(), nil
}
func (e *Ecc) EccSignVer(pt, sign []byte) bool {
	var rInt, sInt big.Int
	// 根据sign，解析出r，s
	rs := bytes.Split(sign, []byte("+"))
	rInt.UnmarshalText(rs[0])
	sInt.UnmarshalText(rs[1])
	// 根据公钥，明文，r，s验证签名
	v := ecdsa.Verify(e.EcdsaPublicKey, pt, &rInt, &sInt)
	return v
}

// GetEcdsaPrivateKey 私钥 -> []byte
func (e *Ecc) GetEcdsaPrivateKey() []byte {
	if e.EcdsaPrivateKey == nil {
		return nil
	}
	privy := e.EcdsaPrivateKey
	return math.PaddedBigBytes(privy.D, privy.Params().BitSize/8)
}

// SetEcdsaPrivateKey []byte -> 私钥
func (e *Ecc) SetEcdsaPrivateKey(d []byte) (err error) {
	e.EcdsaPrivateKey, err = crypto.ToECDSA(d)
	e.EcdsaPublicKey = &e.EcdsaPrivateKey.PublicKey
	e.EciesPrivateKey = ecies.ImportECDSA(e.EcdsaPrivateKey)
	e.EciesPublicKey = &e.EciesPrivateKey.PublicKey
	return err
}

// GetEcdsaPublicKey 公钥 -> []byte
func (e *Ecc) GetEcdsaPublicKey() []byte {
	pub := &e.EcdsaPrivateKey.PublicKey
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(crypto.S256(), pub.X, pub.Y)
}

// SetEcdsaPublicKey []byte -> 公钥
func (e *Ecc) SetEcdsaPublicKey(pub []byte) error {
	if len(pub) == 0 {
		return errors.New("pub len is 0")
	}
	x, y := elliptic.Unmarshal(crypto.S256(), pub)
	e.EcdsaPublicKey = &ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}
	e.EciesPublicKey = ecies.ImportECDSAPublic(e.EcdsaPublicKey)
	return nil
}

func exKey(prv string) *ecies.PrivateKey {
	key, err := crypto.HexToECDSA(prv)
	if err != nil {
		panic(err)
	}
	return ecies.ImportECDSA(key)
}

// SaveECDSA 私钥 -> 文件
// SaveECDSA saves a secp256k1 private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func (e *Ecc) SaveECDSA(file string) error {
	k := hex.EncodeToString(e.GetEcdsaPrivateKey())
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// LoadECDSA 文件 -> 私钥
func (e *Ecc) LoadECDSA(file string) error {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return err
	}

	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return err
	}

	return e.SetEcdsaPrivateKey(key)
}
func TestEcc() {
	var err error
	var SelfEcc = Ecc{}
	var YouEcc = Ecc{}
	var SignPlain = []byte("yemu")
	var SignCipher []byte
	var BytesEcdsaPublicKey []byte
	var Plain = []byte("hello world!")
	var Cipher []byte

	// 生成一对 ECC S256 的 key
	err = SelfEcc.EccGetKey()
	if err != nil {
		panic(err)
	}
	// 导出公钥，并且设置到另一个结构体上
	BytesEcdsaPublicKey = SelfEcc.GetEcdsaPublicKey()
	fmt.Printf("BytesEcdsaPublicKey: %v\n", BytesEcdsaPublicKey)
	err = YouEcc.SetEcdsaPublicKey(BytesEcdsaPublicKey)
	if err != nil {
		panic(err)
	}
	// 签名和验签
	SignCipher, err = SelfEcc.EccSign(SignPlain)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sign: %v\n", YouEcc.EccSignVer(SignPlain, SignCipher))
	// 加密与解密
	Cipher, err = YouEcc.ECCEncrypt(Plain)
	if err != nil {
		panic(err)
	}
	Plain, err = SelfEcc.ECCDecrypt(Cipher)
	fmt.Printf("%s\n", Plain)
	// 私钥保存
	err = SelfEcc.SaveECDSA("ecc.key")
	if err != nil {
		panic(err)
	}
	// 私钥读取
	var NewEcc = Ecc{}
	err = NewEcc.LoadECDSA("ecc.key")
	// 公钥设置
	var New2Ecc = Ecc{}
	err = New2Ecc.SetEcdsaPublicKey([]uint8{4, 136, 76, 169, 50, 219, 86, 251, 42, 210, 193, 174, 161, 229, 226, 177, 94, 177, 86, 1, 224, 132, 80, 145, 168, 124, 130, 66, 176, 5, 140, 186, 73, 19, 226, 205, 234, 10, 44, 65, 8, 108, 205, 64, 2, 157, 63, 5, 79, 184, 110, 225, 197, 187, 78, 255, 27, 83, 169, 209, 3, 146, 211, 130, 230})
	if err != nil {
		panic(err)
	}
	// 新的加解密
	Cipher, err = New2Ecc.ECCEncrypt(Plain)
	if err != nil {
		panic(err)
	}
	Plain, err = NewEcc.ECCDecrypt(Cipher)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", Plain)
	// 新的签名验签
	SignCipher, err = NewEcc.EccSign(SignPlain)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sign: %v\n", New2Ecc.EccSignVer(SignPlain, SignCipher))

}
