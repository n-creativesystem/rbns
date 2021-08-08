package repository

import "github.com/n-creativesystem/rbns/domain/model"

type Permission interface {
	FindAll() (model.Permissions, error)
	FindByID(id model.ID) (*model.Permission, error)
	FindByName(name model.Name) (*model.Permission, error)
}

type PermissionCommand interface {
	Permission
	Create(name model.Name, description string) (*model.Permission, error)
	CreateBatch(names []model.Name, descriptions []string) ([]*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id model.ID) error
}
