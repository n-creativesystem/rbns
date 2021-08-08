package infra

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/infra/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type user struct {
	db *gorm.DB
}

var (
	_ repository.User        = (*user)(nil)
	_ repository.UserCommand = (*user)(nil)
)

func (r *user) FindAll(organizationID model.ID) (model.Users, error) {
	var users []entity.User
	err := r.db.
		Where(&entity.User{OrganizationID: *organizationID.Value()}).
		Preload("Organization").
		Preload("UserRoles.Role.RolePermissions.Permission").Find(&users).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	mUsers := make(model.Users, len(users))
	for idx, user := range users {
		if u, err := user.ConvertModel(); err != nil {
			return nil, err
		} else {
			mUsers[idx] = *u
		}
	}
	return mUsers, nil
}

func (r *user) FindByKey(organizationID model.ID, key model.Key) (*model.User, error) {
	var user entity.User
	err := r.db.
		Where(&entity.User{OrganizationID: *organizationID.Value(), Key: *key.Value()}).
		Preload("Organization").
		Preload("UserRoles.Role.RolePermissions.Permission").Find(&user).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return user.ConvertModel()
}

func (r *user) Create(organizationID model.ID, user *model.User) (*model.User, error) {
	roles := model.Roles(user.GetRole()).Copy()
	userRoles := make([]entity.UserRole, len(user.GetRole()))
	for idx, role := range roles {
		userRoles[idx] = entity.UserRole{
			UserKey:        user.GetKey(),
			RoleID:         *role.GetID(),
			OrganizationID: *organizationID.Value(),
		}
	}
	eUser := entity.User{
		Key:            user.GetKey(),
		OrganizationID: *organizationID.Value(),
	}
	err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&eUser).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	if len(userRoles) > 0 {
		if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error; err != nil {
			return nil, model.NewDBErr(err)
		}
	}
	return user, nil
}

func (r *user) Delete(organizationID model.ID, key model.Key) error {
	db := r.db.
		Where(&entity.User{
			Key:            *key.Value(),
			OrganizationID: *organizationID.Value(),
		}).
		Delete(&entity.User{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}

func (r *user) AddRole(organizationID model.ID, key model.Key, roles ...model.Role) error {
	userRoles := make([]entity.UserRole, len(roles))
	for idx, role := range roles {
		userRoles[idx] = entity.UserRole{
			OrganizationID: *organizationID.Value(),
			UserKey:        *key.Value(),
			RoleID:         *role.GetID(),
		}
	}
	return model.NewDBErr(r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error)
}

func (r *user) DeleteRole(organizationID model.ID, key model.Key, roleId model.ID) error {
	db := r.db.Where(&entity.UserRole{
		UserKey:        *key.Value(),
		OrganizationID: *organizationID.Value(),
		RoleID:         *roleId.Value(),
	}).Delete(&entity.UserRole{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}
