package p100

import "errors"

var (
	ErrNoCookieToken  = errors.New("cookie token is empty, please call the Handshake function")
	ErrNoRequestToken = errors.New("request token is empty, please call the Login function")
	ErrInvalidUrl     = errors.New("url is not set to a valid value")
	ErrNoEncryptor    = errors.New("encryptor is not set, please call Handshake function")
)
