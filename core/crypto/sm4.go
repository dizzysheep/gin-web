package crypto

import (
	"crypto/cipher"
	"errors"

	"github.com/tjfoc/gmsm/sm4"
)

func Sm4Encrypt(key, iv, plantText []byte, paddingStatus bool) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if paddingStatus {
		plantText = pkcs5Padding(plantText, block.BlockSize())
	}

	blockModel := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func Sm4Decrypt(key, iv, ciphertext []byte, paddingStatus bool) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < block.BlockSize() {
		return nil, errors.New("Sm4Decrypt crypto/cipher: ciphertext too short")
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("Sm4Decrypt crypto/cipher: ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	if paddingStatus {
		return pkcs5UnPadding(plantText)
	}

	return plantText, nil
}
