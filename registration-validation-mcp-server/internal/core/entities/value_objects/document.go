package valueobjects

const (
	CPFLength  = 11
	CNPJLength = 14
)

type Document string

func NewDocument(document string) Document {
	return Document(document)
}

func (d Document) String() string {
	return string(d)
}

func ValidateDocument(document string) bool {
	switch len(document) {
	case CPFLength:
		return validateCpf(document)
	case CNPJLength:
		return validateCnpj(document)
	default:
		return false
	}
}

func validateCpf(cpf string) bool {
	validation := func(d string) bool {
		if blacklist(d, CPFLength) {
			return false
		}

		var sum int
		for i := 0; i < 9; i++ {
			sum += int(d[i]-'0') * (10 - i)
		}
		firstCheck := (sum * 10) % 11
		if firstCheck == 10 {
			firstCheck = 0
		}
		if firstCheck != int(d[9]-'0') {
			return false
		}

		sum = 0
		for i := 0; i < 10; i++ {
			sum += int(d[i]-'0') * (11 - i)
		}
		secondCheck := (sum * 10) % 11
		if secondCheck == 10 {
			secondCheck = 0
		}
		return secondCheck == int(d[10]-'0')
	}

	return validation(cpf)
}

func validateCnpj(cnpj string) bool {
	validation := func(d string) bool {
		if blacklist(d, CNPJLength) {
			return false
		}

		var sum int
		weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < 12; i++ {
			sum += int(d[i]-'0') * weights1[i]
		}
		firstCheck := sum % 11
		if firstCheck < 2 {
			firstCheck = 0
		} else {
			firstCheck = 11 - firstCheck
		}
		if firstCheck != int(d[12]-'0') {
			return false
		}

		sum = 0
		weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < 13; i++ {
			sum += int(d[i]-'0') * weights2[i]
		}
		secondCheck := sum % 11
		if secondCheck < 2 {
			secondCheck = 0
		} else {
			secondCheck = 11 - secondCheck
		}
		return secondCheck == int(d[13]-'0')
	}

	return validation(cnpj)
}

// blacklist checks for sequences of the same digit
func blacklist(number string, length int) bool {
	for i := 1; i < length; i++ {
		if number[i] != number[0] {
			return false
		}
	}
	return true
}
