package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strings"
	"time"
)

// aes 加密
type aesCrypto struct {
	// 密钥必须时 16 位
	key []byte

	// 加密模式
	mode string

	// 用于生成初始向量
	t time.Time
}

type aesCbcCrypto struct {
	// 密钥必须时 16 位
	key []byte

	// iv 向量
	iv []byte
}

// 创建 aes 实例
func NewAes(key []byte, mode string) *aesCrypto {
	return &aesCrypto{key, strings.ToLower(mode), time.Now()}
}

// 创建 aes cbc实例
func NewAesCBC(key, vi []byte) *aesCbcCrypto {
	return &aesCbcCrypto{key, vi}
}

// aes 加密
func (a *aesCrypto) Encrypt(data []byte) ([]byte, error) {
	switch a.mode {
	case ECB:
		return a.ecbEncrypt(data)
	case CBC:
		return a.cbcEncrypt(data)
	default:
		return nil, errors.New("aes encrypt not support mode: " + a.mode)
	}
}

// aes 解密
func (a *aesCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	switch a.mode {
	case ECB:
		return a.ecbDecrypt(ciphertext)
	case CBC:
		return a.cbcDecrypt(ciphertext)
	default:
		return nil, errors.New("aes decrypt not support mode: " + a.mode)
	}
}

// aes cbc 模式加密
func (a *aesCrypto) cbcEncrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = pkcs5Padding(data, blockSize)
	iv := []byte(a.t.String())[:blockSize] // 初始向量
	blockMode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(data))
	blockMode.CryptBlocks(dst, data)
	return dst, nil
}

// aes cbc 模式解密
func (a *aesCrypto) cbcDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	iv := []byte(a.t.String())[:block.BlockSize()] // 初始向量
	blockMode := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(out, ciphertext)
	return pkcs5UnPadding(out)
}

// aes ecb 模式加密
func (a *aesCrypto) ecbEncrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}
	dst := make([]byte, len(plain))

	// 分组分块加密
	blockSize := block.BlockSize()
	for bs, be := 0, blockSize; bs <= len(data); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(dst[bs:be], plain[bs:be])
	}
	return dst, nil
}

// aes ecb 模式解密
func (a *aesCrypto) ecbDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(ciphertext))
	blockSize := block.BlockSize()
	for bs, be := 0, blockSize; bs < len(ciphertext); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(out[bs:be], ciphertext[bs:be])
	}
	trim := 0
	if len(out) > 0 {
		trim = len(out) - int(out[len(out)-1])
	}
	return out[:trim], nil
}

// aes cbc 模式加密
func (a *aesCbcCrypto) CbcEncrypt(plantText []byte, padding bool) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	if padding {
		plantText = pkcs5Padding(plantText, block.BlockSize())
	}

	blockMode := cipher.NewCBCEncrypter(block, a.iv)
	ciphertext := make([]byte, len(plantText))
	blockMode.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

// aes cbc 模式解密
func (a *aesCbcCrypto) CbcDecrypt(ciphertext []byte, padding bool) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, a.iv)
	plantText := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(plantText, ciphertext)
	if padding {
		return pkcs5UnPadding(plantText)
	}

	return plantText, nil
}
