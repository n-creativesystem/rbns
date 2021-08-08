package repository

import (
	"github.com/n-creativesystem/rbns/domain/model"
)

type User interface {
	FindAll(organizationID model.ID) (model.Users, error)
	FindByKey(organizationID model.ID, key model.Key) (*model.User, error)
}

type UserCommand interface {
	User
	Create(organizationID model.ID, user *model.User) (*model.User, error)
	Delete(organizationID model.ID, key model.Key) error
	AddRole(organizationID model.ID, key model.Key, roles ...model.Role) error
	DeleteRole(organizationID model.ID, key model.Key, roleId model.ID) error
}
