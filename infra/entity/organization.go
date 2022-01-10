package entity

import (
	"fmt"

	"github.com/n-creativesystem/rbns/domain/model"
)

type Organization struct {
	Model
	Name        string   `gorm:"type:varchar(256);uniqueIndex"`
	Description string   `gorm:"type:varchar(256)"`
	Users       []User   `gorm:"many2many:organization_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Roles       []Role   `gorm:"many2many:organization_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tenants     []Tenant `gorm:"many2many:tenant_organizations;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Organization) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_tanant", table)
}

func (org Organization) ConvertModel() *model.Organization {
	o := &model.Organization{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		Users:       make([]model.User, 0, len(org.Roles)),
		Roles:       make([]model.Role, 0, len(org.Roles)),
	}

	for _, role := range org.Roles {
		o.Roles = append(o.Roles, *role.ConvertModel())
	}

	for _, user := range org.Users {
		o.Users = append(o.Users, *user.ConvertModel())
	}

	return o
}
