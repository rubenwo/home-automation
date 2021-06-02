package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
)

type TapoEncryptor struct {
	iv  []byte
	key []byte
}

func NewTapoEncryptor(key, iv []byte) *TapoEncryptor {
	return &TapoEncryptor{
		iv:  iv,
		key: key,
	}
}

func (t *TapoEncryptor) Encrypt(data []byte) ([]byte, error) {
	plaintext := data

	plaintext, err := pkcs7Padding(plaintext)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher(t.key)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, t.iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func (t *TapoEncryptor) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(t.key)
	if err != nil {
		log.Fatal(err)
	}
	plainText := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, t.iv)
	mode.CryptBlocks(plainText, data)

	plainText, err = pkcs7UnPadding(plainText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func pkcs7Padding(ciphertext []byte) ([]byte, error) {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...), nil
}

func pkcs7UnPadding(plainText []byte) ([]byte, error) {
	length := len(plainText)
	if length == 0 {
		return nil, errors.New("plainText is empty, thus could not un-pad")
	}
	unpadding := int(plainText[length-1])
	if length-unpadding > len(plainText) {
		return nil, errors.New("")
	}
	return plainText[:(length - unpadding)], nil
}
