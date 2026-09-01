package valueobjects

import (
	"regexp"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"golang.org/x/crypto/bcrypt"
)

const PasswordPattern = `^[A-Za-z0-9!@#$%^&*()_+]{8,}$`

type Password string

func NewPassword(password string) Password {
	return Password(password)
}

func ValidatePassword(p string) bool {
	return len(p) >= 8 && regexp.MustCompile(PasswordPattern).MatchString(p)
}

func (p Password) String() string {
	return string(p)
}

func (p Password) Hash() (Password, *rest_err.RestErr) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", rest_err.NewInternalServerError("failed to hash password").WithCause(err)
	}
	return Password(pwd), nil
}

func (p Password) CompareHashWithPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	return err == nil
}
