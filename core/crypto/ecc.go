package crypto

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/x509"
    "encoding/pem"
)

// ecc 加解密
type eccCrypto struct {
    // 公钥
    publicKey []byte

    // 私钥
    privateKey []byte
}

// 创建 ecc 实例，仅支持 pem 格式密钥
func NewEcc(pubKey, priKey []byte) *eccCrypto {
    return &eccCrypto{pubKey, priKey}
}

// ecc 公钥加密
func (e *eccCrypto) EncryptWithPublicKey(data []byte) ([]byte, error) {
    // 解析 ecc 公钥
    block, _ := pem.Decode(e.publicKey)
    if block == nil {
        return nil, ErrInvalidPublicKey
    }
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    pubKey := ImportECDSAPublic(pub.(*ecdsa.PublicKey))

    cipherText, err := Encrypt(rand.Reader, pubKey, data, nil, nil)
    if err != nil {
        return nil, err
    }
    return cipherText, nil
}

// ecc 私钥解密
func (e *eccCrypto) DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
    // 解析 ecc 私钥
    block, _ := pem.Decode(e.privateKey)
    if block == nil {
        return nil, ErrInvalidPrivateKey
    }
    pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    priKey := ImportECDSA(pri.(*ecdsa.PrivateKey))
    dst, err := priKey.Decrypt(ciphertext, nil, nil)
    if err != nil {
        return nil, err
    }
    return dst, nil
}

// ecc 私钥加密 TODO
func (e *eccCrypto) EncryptWithPrivateKey(data []byte) ([]byte, error) {
    return nil, nil
}

// ecc 公钥解密 TODO
func (e *eccCrypto) DecryptWithPublicKey(ciphertext []byte) ([]byte, error) {
    return nil, nil
}
