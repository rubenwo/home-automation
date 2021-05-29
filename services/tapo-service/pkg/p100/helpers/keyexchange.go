package helpers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100/encryption"
	"io/ioutil"
	"net/http"
	"net/url"
)

func InitializeEncryption(u *url.URL) (encryption.Encryptor, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, "", ErrGeneratingKey
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, "", ErrMarshallingPrivateKey
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, "", ErrMarshallingPublicKey
	}
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	buffer := bytes.NewBuffer(nil)
	err = pem.Encode(buffer, privateKeyBlock)
	if err != nil {
		return nil, "", ErrEncodingPrivateKey
	}
	//privateKeyString := string(buffer.Bytes())

	buffer.Reset()
	err = pem.Encode(buffer, publicKeyBlock)
	if err != nil {
		return nil, "", ErrEncodingPublicKey
	}

	publicKeyString := string(buffer.Bytes())

	var body struct {
		Method string `json:"method"`
		Params struct {
			Key string `json:"key"`
		} `json:"params"`
	}
	body.Params.Key = publicKeyString
	body.Method = "handshake"

	client := &http.Client{}
	data, _ := json.Marshal(&body)
	resp, err := client.Post(u.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", ErrSendingHandshakeRequest, err.Error())
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", ErrSendingHandshakeRequest, err.Error())
	}
	var response struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Key string `json:"key"`
		} `json:"result"`
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, "", err
	}
	if response.ErrorCode != 0 {
		return nil, "", ErrSendingHandshakeRequest
	}
	cookieToken, err := findCookie(resp, "TP_SESSIONID")
	if err != nil {
		return nil, "", err
	}

	encryptor, err := decodeHandshakeKey(response.Result.Key, privateKey)
	if err != nil {
		return nil, "", err
	}
	return encryptor, cookieToken, nil
}

func findCookie(resp *http.Response, name string) (string, error) {
	for _, cookie := range resp.Cookies() {
		if name == cookie.Name {
			return cookie.Value, nil
		}
	}
	return "", errors.New("no cookie found")
}
func decodeHandshakeKey(key string, privateKey *rsa.PrivateKey) (encryption.Encryptor, error) {
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	decryptedPKCS1v15, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}

	keyArr := make([]byte, 16)
	ivArr := make([]byte, 16)

	for i := 0; i < 16; i++ {
		keyArr[i] = decryptedPKCS1v15[i]
		ivArr[i] = decryptedPKCS1v15[i+16]
	}

	return encryption.NewTapoEncryptor(keyArr, ivArr), nil
}
