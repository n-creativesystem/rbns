package entity

import (
	"github.com/n-creativesystem/rbns/domain/model"
)

type Role struct {
	Model
	Name          string         `gorm:"type:varchar(256);uniqueIndex"`
	Description   string         `gorm:"type:varchar(256)"`
	Permissions   []Permission   `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users         []User         `gorm:"many2many:role_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tenants       []Tenant       `gorm:"many2many:tenant_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Organizations []Organization `gorm:"many2many:organization_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (r *Role) ConvertModel() *model.Role {
	role := &model.Role{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description,
		Permissions:   make([]model.Permission, 0, len(r.Permissions)),
		Organizations: make([]model.Organization, 0, len(r.Organizations)),
	}

	for _, permission := range r.Permissions {
		role.Permissions = append(role.Permissions, *permission.ConvertModel())
	}

	for _, organization := range r.Organizations {
		role.Organizations = append(role.Organizations, *organization.ConvertModel())
	}

	return role
}
