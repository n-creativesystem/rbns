package infra

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/infra/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type role struct {
	db *gorm.DB
}

var (
	_ repository.Role        = (*role)(nil)
	_ repository.RoleCommand = (*role)(nil)
)

func (r *role) FindAll() (model.Roles, error) {
	session := r.db
	var roles []entity.Role
	err := session.Preload("UserRoles.Organization").Preload("RolePermissions.Permission").Order("id").Find(&roles).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	mRoles := make(model.Roles, len(roles))
	for i, role := range roles {
		if r, err := role.ConvertModel(); err != nil {
			return nil, err
		} else {
			mRoles[i] = *r
		}
	}
	return mRoles, nil
}

func (r *role) FindByID(id model.ID) (*model.Role, error) {
	session := r.db
	var role entity.Role
	err := session.Preload("UserRoles.Organization").Preload("RolePermissions.Permission").Where(&entity.Role{Model: entity.Model{ID: *id.Value()}}).First(&role).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return role.ConvertModel()
}

func (r *role) Create(name model.Name, description string) (*model.Role, error) {
	entity := entity.Role{
		Name:        *name.Value(),
		Description: description,
	}
	entity.Generate()
	err := r.db.Create(&entity).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return model.NewRole(entity.ID, entity.Name, entity.Description, nil)
}

func (r *role) CreateBatch(names []model.Name, descriptions []string) ([]*model.Role, error) {
	entities := make([]entity.Role, len(names))
	for idx, name := range names {
		entity := entity.Role{
			Name:        *name.Value(),
			Description: descriptions[idx],
		}
		entity.Generate()
		entities[idx] = entity
	}
	err := r.db.Create(&entities).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	mRoles := make([]*model.Role, len(entities))
	for idx, entity := range entities {
		mRoles[idx], err = model.NewRole(entity.ID, entity.Name, entity.Description, nil)
		if err != nil {
			return nil, err
		}
	}
	return mRoles, nil
}

func (r *role) Update(role *model.Role) error {
	value := entity.Role{
		Name:        *role.GetName(),
		Description: role.GetDescription(),
	}
	return model.NewDBErr(r.db.Where(&entity.Role{Model: entity.Model{ID: *role.GetID()}}).Updates(&value).Error)
}

func (r *role) Delete(id model.ID) error {
	db := r.db.Where(&entity.Role{Model: entity.Model{ID: *id.Value()}}).Delete(&entity.Role{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}

func (r *role) AddPermission(id model.ID, permissions model.Permissions) error {
	roleAndPermissions := make([]entity.RolePermission, len(permissions))
	for idx, permission := range permissions {
		roleAndPermissions[idx] = entity.RolePermission{
			RoleID:       *id.Value(),
			PermissionID: *permission.GetID(),
		}
	}
	return model.NewDBErr(r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roleAndPermissions).Error)
}

func (r *role) DeletePermission(id model.ID, permissionId model.ID) error {
	rolePermission := entity.RolePermission{
		RoleID:       *id.Value(),
		PermissionID: *permissionId.Value(),
	}
	db := r.db.Where(&rolePermission).Delete(&entity.RolePermission{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}
