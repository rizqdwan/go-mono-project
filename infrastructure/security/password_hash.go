package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash interface {
	Hash(password string) (hashedPassword string, err error)
	Compare(password, hashedPassword string) (match bool, err error)
}

type passwordHash struct{}

func NewPasswordHash() PasswordHash {
	return &passwordHash{}
}

func (p *passwordHash) Hash(password string) (hashedPassword string, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (p *passwordHash) Compare(plainPassword, hashedPassword string) (match bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err != nil {
		return false, err
	}

	match = true
	return match, nil
}
