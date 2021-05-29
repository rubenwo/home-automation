package helpers

import "errors"

var (
	ErrGeneratingKey           = errors.New("could not generate RSA private key")
	ErrMarshallingPrivateKey   = errors.New("could not marshal private key")
	ErrEncodingPrivateKey      = errors.New("could not encode private key")
	ErrMarshallingPublicKey    = errors.New("could not marshall public key")
	ErrEncodingPublicKey       = errors.New("could not encode public key")
	ErrSendingHandshakeRequest = errors.New("could not send handshake request")
)
