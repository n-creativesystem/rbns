package entity

import (
	"github.com/n-creativesystem/rbns/domain/model"
)

type User struct {
	Model
	Name          string         `gorm:"uniqueIndex"`
	Organizations []Organization `gorm:"many2many:organization_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Roles         []Role         `gorm:"many2many:role_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (u *User) ConvertModel() *model.User {
	user := &model.User{
		ID:   u.ID.String(),
		Name: u.Name,

		Organizations: make([]model.Organization, 0, len(u.Organizations)),
		Roles:         make([]model.Role, 0, len(u.Roles)),
	}

	for _, org := range u.Organizations {
		user.Organizations = append(user.Organizations, *org.ConvertModel())
	}

	for _, role := range u.Roles {
		user.Roles = append(user.Roles, *role.ConvertModel())
	}

	return user
}
