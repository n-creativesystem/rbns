package res

import "github.com/n-creativesystem/rbns/domain/model"

type User struct {
	Key            string `json:"key"`
	OrganizationId string `json:"organization_id"`
	Roles
	Permissions
}

func NewUser(user model.User) User {
	return User{
		Key:         user.GetKey(),
		Roles:       NewRoles(user.GetRole()),
		Permissions: NewPermissions(user.GetPermission()),
	}
}

func NewUsers(users model.Users) []User {
	us := make([]User, len(users))
	for idx, u := range users {
		us[idx] = NewUser(u)
	}
	return us
}
