package entity

import (
	"time"

	"github.com/n-creativesystem/rbns/domain/model"
)

type Resource struct {
	ID          string `gorm:"type:varchar(256);primaryKey"`
	Description string `gorm:"type:varchar(256)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	ResourcePermissions []ResourcePermissions `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ResourcePermissions struct {
	ResourceID   string      `gorm:"type:varchar(256);primaryKey"`
	PermissionID string      `gorm:"type:varchar(256);primaryKey"`
	Permission   *Permission `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (r *Resource) ConvertModel() (*model.Resource, error) {
	permissions := make(model.Permissions, 0, len(r.ResourcePermissions))
	for _, permission := range r.ResourcePermissions {
		if permission.Permission != nil {
			if p, err := permission.Permission.ConvertModel(); err != nil {
				return nil, err
			} else {
				permissions = append(permissions, *p)
			}
		}
	}
	return &model.Resource{
		ID:          r.ID,
		Description: r.Description,
		Permissions: permissions,
	}, nil
}
