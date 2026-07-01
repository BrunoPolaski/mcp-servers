package valueobjects

const EmailPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

type Email string

func NewEmail(email string) Email {
	return Email(email)
}

func (e Email) String() string {
	return string(e)
}
