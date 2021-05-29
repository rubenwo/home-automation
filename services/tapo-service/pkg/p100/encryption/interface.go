package encryption

type Encryptor interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) []byte
}
