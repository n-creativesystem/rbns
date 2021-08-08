package entity

import "github.com/n-creativesystem/rbns/domain/model"

type Permission struct {
	Model
	Name        string `gorm:"type:varchar(256);UNIQUE"`
	Description string `gorm:"type:varchar(256)"`
	// Roles       []RolePermission `gorm:"foreignKey:PermissionID"`
}

func (p Permission) ConvertModel() (*model.Permission, error) {
	return model.NewPermission(p.ID, p.Name, p.Description)
}
