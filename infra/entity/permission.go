package entity

import (
	"fmt"

	"github.com/n-creativesystem/rbns/domain/model"
)

type Permission struct {
	Model
	Name        string `gorm:"type:varchar(256);uniqueIndex"`
	Description string `gorm:"type:varchar(256)"`
	Roles       []Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Tenants []Tenant `gorm:"many2many:tenant_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Permission) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_tanant", table)
}

func (p Permission) ConvertModel() *model.Permission {
	permission := &model.Permission{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Roles:       make([]model.Role, 0, len(p.Roles)),
	}

	for _, role := range p.Roles {
		permission.Roles = append(permission.Roles, *role.ConvertModel())
	}

	return permission
}
