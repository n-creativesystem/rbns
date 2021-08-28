package entity

import (
	"github.com/n-creativesystem/rbns/domain/model"
)

type Role struct {
	Model
	Name            string           `gorm:"type:varchar(256);UNIQUE"`
	Description     string           `gorm:"type:varchar(256)"`
	RolePermissions []RolePermission `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserRoles       []UserRole       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type RolePermission struct {
	RoleID       string      `gorm:"type:varchar(256);primaryKey"`
	PermissionID string      `gorm:"type:varchar(256);primaryKey"`
	Permission   *Permission `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (r Role) ConvertModel() (*model.Role, error) {
	permissions := model.Permissions{}
	for _, permission := range r.RolePermissions {
		if permission.Permission != nil {
			if p, err := permission.Permission.ConvertModel(); err != nil {
				return nil, err
			} else {
				permissions = append(permissions, *p)
			}
		}
	}
	organizationUserRoles := make(model.OrganizationUserRoles, len(r.UserRoles))
	for idx, userRole := range r.UserRoles {
		org, err := userRole.Organization.ConvertModel()
		if err != nil {
			return nil, err
		}
		orgUserRole, err := model.NewOrganizationUserRole(*org, userRole.UserKey, userRole.RoleID)
		if err != nil {
			return nil, err
		}
		organizationUserRoles[idx] = *orgUserRole
	}
	return model.NewRole(r.ID, r.Name, r.Description, permissions, organizationUserRoles...)
}
