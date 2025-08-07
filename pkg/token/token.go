package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"plaintext"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG. This will return an error if
	// the CSPRNG fails to function correctly.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our
	// tokens, so we use the WithPadding(base32.NoPadding) method in the line
	// below to omit them
	token.PlainText = base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(randomBytes)

	// sha256.Sum256() function returns an *array* of length 32, so to make it
	// easier to work with we convert it to a slice using the [:] operator before
	// storing it.
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
