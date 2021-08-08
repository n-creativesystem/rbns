package repository

import "github.com/n-creativesystem/rbns/domain/model"

type Role interface {
	FindAll() (model.Roles, error)
	FindByID(id model.ID) (*model.Role, error)
}

type RoleCommand interface {
	Role
	Create(name model.Name, description string) (*model.Role, error)
	CreateBatch(names []model.Name, descriptions []string) ([]*model.Role, error)
	Update(role *model.Role) error
	Delete(id model.ID) error
	AddPermission(id model.ID, permissions model.Permissions) error
	DeletePermission(id model.ID, permissionId model.ID) error
}
