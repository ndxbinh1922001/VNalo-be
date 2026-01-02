package value_object

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
)

// Email value object - immutable
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new email value object with validation
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}

// Equals compares two email value objects
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

