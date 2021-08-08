package infra

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/dao/driver/postgres"
	"github.com/n-creativesystem/rbns/infra/entity"
	"gorm.io/gorm"
)

type resource struct {
	db *gorm.DB
}

func (r *resource) Save(method, uri string, permissions ...model.Permission) error {
	for _, p := range permissions {
		if id := p.GetID(); id != nil {
			strId := *id
			resource := entity.Resource{
				Method:       method,
				URI:          uri,
				PermissionID: strId,
			}
			if err := r.db.Create(resource).Error; err != nil {
				if postgres.IsDuplication(err) {
					return postgres.NewDBErr(err)
				}
			}
		}
	}
	return nil
}

func (r *resource) Find(method, uri string) *model.Resource {
	var resources []*entity.Resource
	err := r.db.Where(&entity.Resource{
		Method: method,
		URI:    uri,
	}).Preload("Permission").Find(resources).Error
	if err != nil {
		return nil
	}
	permissions := make([]string, len(resources))
	for i, resource := range resources {
		permissions[i] = resource.Permission.Name
	}
	return model.NewResource(method, uri, permissions...)
}

func (r *resource) FindByMethod(method string) *model.Resource {
	var resources []*entity.Resource
	err := r.db.Where(&entity.Resource{
		Method: method,
	}).Preload("Permission").Find(resources).Error
	if err != nil {
		return nil
	}
	uri := make([]string, len(resources))
	permissions := make([]string, len(resources))
	for i, resource := range resources {
		uri[i] = resource.URI
		permissions[i] = resource.Permission.Name
	}
	return model.NewResourceByMethod(method, uri, permissions)
}
