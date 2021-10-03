package entity

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"gorm.io/gorm"
)

type Organization struct {
	Model
	Name        string     `gorm:"type:varchar(256);UNIQUE"`
	Description string     `gorm:"type:varchar(256)"`
	Users       []User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserRoles   []UserRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (org Organization) ConvertModel() (*model.Organization, error) {
	users := make(model.Users, len(org.Users))
	for idx, user := range org.Users {
		mpPermissions := map[string]model.Permission{}
		roles := make(model.Roles, len(user.UserRoles))
		for j, role := range user.UserRoles {
			r, err := model.NewRole(role.RoleID, role.Role.Name, role.Role.Description, nil)
			if err != nil {
				return nil, err
			}
			roles[j] = *r
			for _, rp := range role.Role.RolePermissions {
				if _, ok := mpPermissions[rp.PermissionID]; !ok {
					p := *rp.Permission
					mP, err := model.NewPermission(p.ID, p.Name, p.Description)
					if err != nil {
						return nil, err
					}
					mpPermissions[rp.PermissionID] = *mP
				}
			}
		}
		j := 0
		permissions := make(model.Permissions, len(mpPermissions))
		for _, val := range mpPermissions {
			permissions[j] = val
			j++
		}
		u, err := model.NewUser(user.Key, roles, permissions)
		if err != nil {
			return nil, err
		}
		users[idx] = *u
	}
	return model.NewOrganization(org.ID, org.Name, org.Description, users...)
}

func organizationMigration(db *gorm.DB) error {
	return db.AutoMigrate(&Organization{}, &User{}, &UserRole{})
}

func organizationMigrationBack(db *gorm.DB) error {
	return db.Migrator().DropTable(&Organization{}, &User{}, &UserRole{})
}
