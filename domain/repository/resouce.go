package repository

import "github.com/n-creativesystem/rbns/domain/model"

type Resource interface {
	FindByMethod(method string) *model.Resource
	Find(method, uri string) *model.Resource
}

type ResourceCommand interface {
	Save(method, uri string, permissions ...model.Permission) error
}
