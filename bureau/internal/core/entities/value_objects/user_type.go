package valueobjects

type UserType string

const (
	UserTypePerson  UserType = "person"
	UserTypeAdmin   UserType = "admin"
	UserTypeAnalyst UserType = "analyst"
)

var mapUserType = map[UserType]struct{}{
	UserTypePerson:  {},
	UserTypeAdmin:   {},
	UserTypeAnalyst: {},
}

func NewUserType(userType string) UserType {
	return UserType(userType)
}

func ValidateUserType(userType string) bool {
	_, ok := mapUserType[UserType(userType)]
	return ok
}

func (dt UserType) String() string {
	return string(dt)
}

func (dt UserType) GreaterOrEqualThan(other UserType) bool {
	order := map[UserType]int{
		UserTypePerson:  1,
		UserTypeAnalyst: 2,
		UserTypeAdmin:   3,
	}

	return order[dt] >= order[other]
}
