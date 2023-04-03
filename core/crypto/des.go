package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"strings"
	"time"
)

const (
	// des 加密模式
	CBC string = "cbc"
	ECB string = "ecb"
)

// des 加解密
type desCrypto struct {
	// 密钥必须是 8 个字节
	key []byte

	// 加密模式
	mode string

	// 用于生成初始向量
	t time.Time
}

// 创建 des 实例
func NewDes(key []byte, mode string) *desCrypto {
	return &desCrypto{key, strings.ToLower(mode), time.Now()}
}

// des 加密
func (d *desCrypto) Encrypt(data []byte) ([]byte, error) {
	switch d.mode {
	case ECB:
		return d.ecbEncrypt(data)
	case CBC:
		return d.cbcEncrypt(data)
	default:
		return nil, errors.New("des encrypt not support mode: " + d.mode)
	}
}

// des 解密
func (d *desCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	switch d.mode {
	case ECB:
		return d.ecbDecrypt(ciphertext)
	case CBC:
		return d.cbcDecrypt(ciphertext)
	default:
		return nil, errors.New("des decrypt not support mode: " + d.mode)
	}
}

// ecb 模式加密
func (d *desCrypto) ecbEncrypt(data []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}

	// 加密数据
	bs := block.BlockSize()
	data = pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		// 对明文按照 blocksize 进行分块加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// ecb 模式解密
func (d *desCrypto) ecbDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(ciphertext))
	dst := out
	for len(ciphertext) > 0 {
		block.Decrypt(dst, ciphertext[:bs])
		ciphertext = ciphertext[bs:]
		dst = dst[bs:]
	}

	return pkcs5UnPadding(out)
}

// cbc 模式加密，安全性高于 ecb
func (d *desCrypto) cbcEncrypt(data []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	data = pkcs5Padding(data, block.BlockSize())
	iv := []byte(d.t.String())[:block.BlockSize()] // 初始向量
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	return out, nil
}

func (d *desCrypto) cbcDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	iv := []byte(d.t.String())[:block.BlockSize()] // 初始向量
	mode := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(ciphertext))
	mode.CryptBlocks(out, ciphertext)

	return pkcs5UnPadding(out)
}

// 3DES 加解密，相当于 des 加解密重复 3 次
// 增加复杂度，以提高破解难度
type tripleDesCrypto struct {
	// 密钥必须是 24 个字节
	key []byte

	// 加密模式
	mode string

	// 用于生成初始向量
	t time.Time
}

// 创建 3des 实例
func New3Des(key []byte, mode string) *tripleDesCrypto {
	return &tripleDesCrypto{key, mode, time.Now()}
}

// 3des 加密
func (d *tripleDesCrypto) TripleEncrypt(data []byte) ([]byte, error) {
	switch strings.ToLower(d.mode) {
	case ECB:
		return d.tripleEcbEncrypt(data)
	case CBC:
		return d.tripleCbcEncrypt(data)
	default:
		return nil, errors.New("3des encrypt not support mode: " + d.mode)
	}
}

// 3des 解密
func (d *tripleDesCrypto) TripleDecrypt(ciphertext []byte) ([]byte, error) {
	switch strings.ToLower(d.mode) {
	case ECB:
		return d.tripleEcbDecrypt(ciphertext)
	case CBC:
		return d.tripleCbcDecrypt(ciphertext)
	default:
		return nil, errors.New("3des decrypt not support mode: " + d.mode)
	}
}

// 3des ecb 加密
func (d *tripleDesCrypto) tripleEcbEncrypt(data []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}

	// 加密数据
	bs := block.BlockSize()
	data = pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		// 对明文按照 blocksize 进行分块加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// 3des ecb 解密
func (d *tripleDesCrypto) tripleEcbDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(ciphertext))
	dst := out
	for len(ciphertext) > 0 {
		block.Decrypt(dst, ciphertext[:bs])
		ciphertext = ciphertext[bs:]
		dst = dst[bs:]
	}
	return pkcs5UnPadding(out)
}

// 3des cbc 加密
func (d *tripleDesCrypto) tripleCbcEncrypt(data []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}
	data = pkcs5Padding(data, block.BlockSize())
	iv := []byte(d.t.String())[:block.BlockSize()] // 初始向量
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	return out, nil
}

// 3des cbc 解密
func (d *tripleDesCrypto) tripleCbcDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}
	iv := []byte(d.t.String())[:block.BlockSize()] // 初始向量
	mode := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(ciphertext))
	mode.CryptBlocks(out, ciphertext)

	return pkcs5UnPadding(out)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs5UnPadding(plantText []byte) ([]byte, error) {
	length := len(plantText)
	unpadding := int(plantText[length-1])

	len := length - unpadding
	if len < 0 || len > length {
		return nil, errors.New("aes unpadding error")
	}

	return plantText[:(len)], nil
}
