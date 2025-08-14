package user

import (
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	"github.com/barantoraman/microgate/pkg/validator"
)

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid an email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, u *entity.User) {
	ValidateEmail(v, u.Email)
	ValidatePassword(v, u.Password)

	v.Check(u.PasswordHash != nil, "password_hash", "password hash must be provided")
}
