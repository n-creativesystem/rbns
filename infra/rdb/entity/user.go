package entity

import "github.com/n-creativesystem/rbns/domain/model"

type User struct {
	Key            string       `gorm:"type:varchar(256);primaryKey"`
	OrganizationID string       `gorm:"type:varchar(256);primaryKey"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserRoles      []UserRole   `gorm:"foreignKey:UserKey,OrganizationID;references:Key,OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u User) ConvertModel() (*model.User, error) {
	var permissions []model.Permission
	var roles []model.Role
	for _, role := range u.UserRoles {
		if r, err := role.Role.ConvertModel(); err != nil {
			return nil, err
		} else {
			roles = append(roles, *r)
			permissions = append(permissions, r.GetPermissions().Copy()...)
		}

	}
	return model.NewUser(u.Key, roles, permissions)
}

type UserRole struct {
	UserKey        string `gorm:"type:varchar(256);primaryKey"`
	OrganizationID string `gorm:"type:varchar(256);primaryKey"`
	RoleID         string `gorm:"type:varchar(256);primaryKey"`
	Organization   Organization
	Role           Role
}
