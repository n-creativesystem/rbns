package service

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
)

func init() {
	di.MustRegister(newResource)
}

type Resource interface {
	Save(ctx context.Context, method, uri string, permissions ...string) error
	Authorized(ctx context.Context, method, uri, organizationName, userKey string) bool
}

type resource struct {
	repo repository.Repository
}

func newResource(repo repository.Repository) Resource {
	return &resource{
		repo: repo,
	}
}

func (r *resource) Save(ctx context.Context, method, uri string, permissions ...string) error {
	con := r.repo.NewConnection()
	permissionsRepo := con.Permission(ctx)
	var mPermissions []model.Permission
	for _, p := range permissions {
		name, err := model.NewName(p)
		if err != nil {
			return err
		}
		permission, err := permissionsRepo.FindByName(name)
		if err != nil {
			return err
		}
		mPermissions = append(mPermissions, *permission)
	}
	return con.Transaction(ctx).Do(func(tx repository.Transaction) error {
		return tx.Resource().Save(method, uri, mPermissions...)
	})
}

func (r *resource) Authorized(ctx context.Context, method, uri, organizationName, userKey string) bool {
	con := r.repo.NewConnection()
	rRepo := con.Resource(ctx)
	uRepo := con.User(ctx)
	oRepo := con.Organization(ctx)
	resources := rRepo.Find(method, uri)
	if resources == nil {
		return false
	}
	if organizationName == "" {
		organizationName = "default"
	}
	oName, err := model.NewName(organizationName)
	if err != nil {
		return false
	}
	organization, err := oRepo.FindByName(oName)
	if err != nil {
		return false
	}
	uKey, err := model.NewKey(userKey)
	if err != nil {
		return false
	}
	oId, _ := model.NewID(*organization.GetID())
	user, err := uRepo.FindByKey(oId, uKey)
	if err != nil {
		return false
	}
	return resources.Check(method, uri, user)
}
