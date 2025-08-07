package token

import "time"

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
}
