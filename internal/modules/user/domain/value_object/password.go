package value_object

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword  = errors.New("password must be at least 8 characters")
	ErrPasswordMismatch = errors.New("password mismatch")
)

// Password value object - immutable and encrypted
type Password struct {
	hashedValue string
}

// NewPassword creates a new password value object with hashing
func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, ErrInvalidPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hashedValue: string(hashedBytes)}, nil
}

// NewPasswordFromHash creates password from already hashed value
func NewPasswordFromHash(hashedPassword string) Password {
	return Password{hashedValue: hashedPassword}
}

// Hash returns the hashed password
func (p Password) Hash() string {
	return p.hashedValue
}

// Compare checks if the plain password matches the hashed password
func (p Password) Compare(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
}

