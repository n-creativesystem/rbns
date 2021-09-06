package rdb

import (
	"errors"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type resource struct {
	db *gorm.DB
}

func (r *resource) Save(id, description string, permissions ...model.Permission) error {
	session := r.db.Session(&gorm.Session{})
	if err := session.Save(&entity.Resource{ID: id, Description: description}).Error; err != nil {
		return err
	}
	resourceId := id
	for _, p := range permissions {
		if id := p.GetID(); id != nil {
			strId := *id
			rp := &entity.ResourcePermissions{
				ResourceID:   resourceId,
				PermissionID: strId,
			}
			if err := session.Save(rp).Error; err != nil {
				return IsDuplication(r.db, err)
			}
		}
	}
	return nil
}

func (r *resource) Delete(id string) error {
	db := r.db.Where(&entity.Resource{ID: id}).Delete(&entity.Resource{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}

func (r *resource) Find(id string) (*model.Resource, error) {
	session := r.db.Session(&gorm.Session{})
	resource := entity.Resource{
		ID: id,
	}
	if err := session.Preload("ResourcePermissions.Permission").First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNoData
		} else {
			logrus.Error(err)
			return nil, err
		}
	}
	return resource.ConvertModel()
}

func (r *resource) FindAll() ([]*model.Resource, error) {
	session := r.db.Session(&gorm.Session{})
	resources := []*entity.Resource{}
	if err := session.Preload("ResourcePermissions.Permission").Find(&resources).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	results := make([]*model.Resource, 0, len(resources))
	for _, resource := range resources {
		result, err := resource.ConvertModel()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *resource) AddPermission(id string, permissions model.Permissions) error {
	resourceAndPermissions := make([]entity.ResourcePermissions, len(permissions))
	for idx, permission := range permissions {
		resourceAndPermissions[idx] = entity.ResourcePermissions{
			ResourceID:   id,
			PermissionID: *permission.GetID(),
		}
	}
	return model.NewDBErr(r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&resourceAndPermissions).Error)
}

func (r *resource) DeletePermission(id string, permissionId model.ID) error {
	resourcePermission := entity.ResourcePermissions{
		ResourceID:   id,
		PermissionID: *permissionId.Value(),
	}
	db := r.db.Where(&resourcePermission).Delete(&entity.ResourcePermissions{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}
