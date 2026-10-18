package models

type Credentials struct {
	Address    string
	Username   string
	Password   string
	PublicKey  []byte
	PrivateKey []byte
}
