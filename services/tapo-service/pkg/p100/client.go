package p100

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100/encryption"
	"github.com/rubenwo/home-automation/services/tapo-service/pkg/p100/helpers"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	URL *url.URL

	encryptor encryption.Encryptor

	cookieToken  string
	requestToken string

	username string
	password string
}

func New(host, username, password string) (*Client, error) {
	u := ""
	if !strings.HasPrefix(host, "http") {
		u = fmt.Sprintf("http://%s/app", host)
	} else {
		u = fmt.Sprintf("%s/app", host)
	}

	uri, err := url.Parse(u)
	if err != nil {
		return nil, ErrInvalidUrl
	}

	c := &Client{
		URL:      uri,
		username: username,
		password: password,
	}

	if err := c.Handshake(); err != nil {
		return nil, err
	}
	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Handshake() error {
	// Reset the tokens because we'll exchange new keys
	c.cookieToken = ""
	c.requestToken = ""
	c.encryptor = nil

	encryptor, cookieToken, err := helpers.InitializeEncryption(c.URL)
	if err != nil {
		return err
	}
	c.encryptor = encryptor
	c.cookieToken = cookieToken

	return nil
}

func shaDigestUsername(username string) string {
	h := sha1.New()
	io.WriteString(h, username)
	digest := h.Sum(nil)

	return hex.EncodeToString(digest)
}

func (c *Client) Login(username, password string) error {
	digestUsername := shaDigestUsername(username)

	var body struct {
		Method string `json:"method"`
		Params struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"params"`
		RequestTimeMils int `json:"requestTimeMils"`
	}

	body.Method = helpers.LoginDevice
	body.Params.Username = base64.StdEncoding.EncodeToString([]byte(digestUsername))
	body.Params.Password = base64.StdEncoding.EncodeToString([]byte(password))
	body.RequestTimeMils = 0

	data, err := json.Marshal(&body)
	if err != nil {
		return err
	}
	encrypted, err := c.encryptor.Encrypt(data)
	if err != nil {
		return err
	}
	var secureBody struct {
		Method string `json:"method"`
		Params struct {
			Request string `json:"request"`
		} `json:"params"`
	}
	secureBody.Method = helpers.SecurePassthrough
	secureBody.Params.Request = string(encrypted)

	data, err = json.Marshal(&secureBody)

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, c.URL.String(), bytes.NewBuffer(data))
	req.AddCookie(&http.Cookie{
		Name:  "TP_SESSIONID",
		Value: c.cookieToken,
	})

	resp, _ := client.Do(req)
	raw, _ := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		return err
	}

	var response struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Response string `json:"response"`
		} `json:"result"`
	}

	err = json.Unmarshal(raw, &response)
	if err != nil {
		return err
	}
	decodeString, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		return err
	}
	decryptedInnerResponse, err := c.encryptor.Decrypt(decodeString)
	if err != nil {
		if err := c.Handshake(); err != nil {
			return err
		}
		if err := c.Login(c.username, c.password); err != nil {
			return err
		}
	}

	var innerResponse struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Token string `json:"token"`
		} `json:"result"`
	}
	err = json.Unmarshal(decryptedInnerResponse, &innerResponse)
	if err != nil {
		return err
	}

	c.requestToken = innerResponse.Result.Token
	return nil
}

func (c *Client) DeviceInfo() (map[string]interface{}, error) {
	if err := c.checkTokenState(); err != nil {
		return nil, err
	}

	var body struct {
		Method string   `json:"method"`
		Params struct{} `json:"params"`
	}
	body.Method = helpers.GetDeviceInfo
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	encrypted, err := c.encryptor.Encrypt(data)
	if err != nil {
		return nil, err
	}
	var secureBody struct {
		Method string `json:"method"`
		Params struct {
			Request string `json:"request"`
		} `json:"params"`
	}
	secureBody.Method = helpers.SecurePassthrough
	secureBody.Params.Request = string(encrypted)

	data, err = json.Marshal(&secureBody)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?token=%s", c.URL.String(), c.requestToken), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Name:  "TP_SESSIONID",
		Value: c.cookieToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	var response struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Response string `json:"response"`
		} `json:"result"`
	}
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return nil, err
	}

	decodeString, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		return nil, err
	}
	decryptedInnerResponse, err := c.encryptor.Decrypt(decodeString)
	if err != nil {
		if err := c.Handshake(); err != nil {
			return nil, err
		}
		if err := c.Login(c.username, c.password); err != nil {
			return nil, err
		}
		if _, err := c.DeviceInfo(); err != nil {
			return nil, err
		}
	}
	var innerResponse struct {
		ErrorCode int                    `json:"error_code"`
		Result    map[string]interface{} `json:"result"`
	}
	err = json.Unmarshal(decryptedInnerResponse, &innerResponse)
	if err != nil {
		return nil, err
	}
	return innerResponse.Result, nil
}

func (c *Client) SetState(deviceOn bool, brightness int) error {
	if err := c.checkTokenState(); err != nil {
		return err
	}
	var body struct {
		Method string `json:"method"`
		Params struct {
			DeviceOn   bool `json:"device_on"`
			Brightness int  `json:"brightness"`
		} `json:"params"`
	}
	body.Method = helpers.SetDeviceInfo
	body.Params.DeviceOn = deviceOn
	body.Params.Brightness = brightness

	data, err := json.Marshal(&body)
	if err != nil {
		return err
	}
	encrypted, err := c.encryptor.Encrypt(data)
	if err != nil {
		return err
	}
	var secureBody struct {
		Method string `json:"method"`
		Params struct {
			Request string `json:"request"`
		} `json:"params"`
	}
	secureBody.Method = helpers.SecurePassthrough
	secureBody.Params.Request = string(encrypted)

	data, err = json.Marshal(&secureBody)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?token=%s", c.URL.String(), c.requestToken), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:  "TP_SESSIONID",
		Value: c.cookieToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	var response struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Response string `json:"response"`
		} `json:"result"`
	}
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return err
	}

	decodeString, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		return err
	}
	decryptedInnerResponse, err := c.encryptor.Decrypt(decodeString)
	if err != nil {
		if err := c.Handshake(); err != nil {
			return err
		}
		if err := c.Login(c.username, c.password); err != nil {
			return err
		}
		if err := c.SetState(deviceOn, brightness); err != nil {
			return err
		}
	}
	var innerResponse struct {
		ErrorCode int                    `json:"error_code"`
		Result    map[string]interface{} `json:"result"`
	}
	err = json.Unmarshal(decryptedInnerResponse, &innerResponse)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Name() string {
	info, err := c.DeviceInfo()
	if err != nil {
		return ""
	}
	name, _ := base64.StdEncoding.DecodeString(info["nickname"].(string))
	return string(name)
}

func (c *Client) checkTokenState() error {
	if c.cookieToken == "" {
		return ErrNoCookieToken
	}
	if c.requestToken == "" {
		return ErrNoRequestToken
	}

	if c.encryptor == nil {
		return ErrNoEncryptor
	}
	return nil
}
