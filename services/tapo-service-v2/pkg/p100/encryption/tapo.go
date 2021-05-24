package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

func (t *TapoEncryptor) Encrypt(data []byte) []byte {
	plaintext := data

	plaintext = pkcs7Padding(plaintext)
	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher(t.key)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, t.iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext))
}

func (t *TapoEncryptor) Decrypt(data []byte) []byte {
	block, err := aes.NewCipher(t.key)
	if err != nil {
		log.Fatal(err)
	}
	plainText := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, t.iv)
	mode.CryptBlocks(plainText, []byte(data))

	plainText = pkcs7UnPadding(plainText)

	return plainText
}

func pkcs7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
