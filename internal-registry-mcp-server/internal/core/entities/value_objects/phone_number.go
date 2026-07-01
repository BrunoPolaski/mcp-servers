package valueobjects

import (
	"regexp"
	"strings"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

const PhonePattern = `^\d{10,11}$`

type PhoneNumber string

func NewPhoneNumber(phone string) (PhoneNumber, *rest_err.RestErr) {
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	phone = strings.TrimLeft(phone, "5")

	if !ValidatePhoneNumber(phone) {
		return "", rest_err.NewBadRequestError("invalid phone number")
	}

	return PhoneNumber(phone), nil
}

func (p PhoneNumber) String() string {
	return string(p)
}

func ValidatePhoneNumber(p string) bool {
	re := regexp.MustCompile(PhonePattern)
	return re.MatchString(p)
}
