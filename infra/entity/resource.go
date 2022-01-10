package entity

// type Resource struct {
// 	ID          plugins.ID     `gorm:"type:varchar(256);primaryKey"`
// 	Tenant      plugins.Tenant `gorm:"type:varchar(256);primaryKey;"`
// 	Description string         `gorm:"type:varchar(256)"`
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time

// 	ResourcePermissions []ResourcePermissions `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
// }

// type ResourcePermissions struct {
// 	ResourceID       string         `gorm:"type:varchar(256);primaryKey"`
// 	PermissionID     string         `gorm:"type:varchar(256);primaryKey"`
// 	PermissionTenant plugins.Tenant `gorm:"type:varchar(256);primaryKey"`
// 	ResourceTenant   plugins.Tenant `gorm:"type:varchar(256);primaryKey"`
// 	Permission       *Permission    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }

// func (r *Resource) ConvertModel() (*model.Resource, error) {
// 	permissions := make(model.Permissions, 0, len(r.ResourcePermissions))
// 	for _, permission := range r.ResourcePermissions {
// 		if permission.Permission != nil {
// 			if p, err := permission.Permission.ConvertModel(); err != nil {
// 				return nil, err
// 			} else {
// 				permissions = append(permissions, *p)
// 			}
// 		}
// 	}
// 	return &model.Resource{
// 		ID:          r.ID.String(),
// 		Description: r.Description,
// 		Permissions: permissions,
// 	}, nil
// }
