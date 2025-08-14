package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash []byte `json:"-"`
}

func (u *User) Set(plaintextPass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPass), 12)
	if err != nil {
		return err
	}
	u.Password = plaintextPass
	u.PasswordHash = hash
	return nil
}

func (u *User) Matches(plaintextPass string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(plaintextPass)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
