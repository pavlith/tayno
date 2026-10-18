package models

type UserOption func(*User)

func WithUsername(username string) UserOption {
	return func(u *User) {
		u.Username = username
	}
}

func WithPasswordHash(hash string) UserOption {
	return func(u *User) {
		u.PasswordHash = hash
	}
}

func WithPublicKey(publicKey []byte) UserOption {
	return func(u *User) {
		u.PublicKey = publicKey
	}
}
