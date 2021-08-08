package service

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
)

func init() {
	di.MustRegister(newPermissionService)
}

type PermissionService interface {
	Create(ctx context.Context, names, descriptions []string) (model.Permissions, error)
	FindById(ctx context.Context, strId string) (*model.Permission, error)
	FindAll(ctx context.Context) (model.Permissions, error)
	Update(ctx context.Context, strId, name, description string) error
	Delete(ctx context.Context, strId string) error
	Check(ctx context.Context, userKey, organizationName string, permissionNames ...string) (*model.ResourceCheck, error)
}
type permissionService struct {
	repo repository.Repository
}

var _ PermissionService = (*permissionService)(nil)

func newPermissionService(repo repository.Repository) PermissionService {
	return &permissionService{repo: repo}
}

// Permission
func (srv *permissionService) Create(ctx context.Context, names, descriptions []string) (model.Permissions, error) {
	var out model.Permissions
	mNames := make([]model.Name, len(names))
	mDescriptions := make([]string, len(descriptions))
	copy(mDescriptions, descriptions)
	for idx, name := range names {
		var err error
		mNames[idx], err = model.NewName(name)
		if err != nil {
			return nil, err
		}
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err := tx.Do(func(tx repository.Transaction) error {
		permissionRepo := tx.Permission()
		permissions, err := permissionRepo.CreateBatch(mNames, mDescriptions)
		if err != nil {
			return err
		}
		out = make(model.Permissions, len(permissions))
		for idx, p := range permissions {
			out[idx] = *p
		}
		return nil
	})
	return out, err
}

func (srv *permissionService) FindById(ctx context.Context, strId string) (*model.Permission, error) {
	permissionRepo := srv.repo.NewConnection().Permission(ctx)
	id, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	permission, err := permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (srv *permissionService) FindAll(ctx context.Context) (model.Permissions, error) {
	permissionRepo := srv.repo.NewConnection().Permission(ctx)
	permissions, err := permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	out := make(model.Permissions, len(permissions))
	copy(out, permissions)
	return out, nil
}

func (srv *permissionService) Update(ctx context.Context, strId, name, description string) error {
	p, err := model.NewPermission(strId, name, description)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Permission().Update(p) }); err != nil {
		return err
	}
	return nil
}

func (srv *permissionService) Delete(ctx context.Context, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Permission().Delete(id) }); err != nil {
		return err
	}
	return nil
}

func (srv *permissionService) Check(ctx context.Context, userKey, organizationName string, permissionNames ...string) (*model.ResourceCheck, error) {
	con := srv.repo.NewConnection()
	mOrganizationName, err := model.NewName(organizationName)
	if err != nil {
		return model.NewResourceCheck(false, err.Error()), err
	}
	mUserKey, err := model.NewKey(userKey)
	if err != nil {
		return model.NewResourceCheck(false, err.Error()), err
	}
	org, err := con.Organization(ctx).FindByName(mOrganizationName)
	if err != nil {
		return model.NewResourceCheck(false, err.Error()), err
	}
	if u, ok := org.IsContainsUsers(mUserKey); !ok {
		return model.NewResourceCheck(false, model.ErrNoData.Error()), model.ErrNoData
	} else {
		for _, permissionName := range permissionNames {
			mPermissionName, err := model.NewName(permissionName)
			if err != nil {
				return model.NewResourceCheck(false, err.Error()), err
			}
			if u.IsContainsPermissionByName(mPermissionName) {
				return model.NewResourceCheck(true, ""), nil
			}
		}
	}
	return model.NewResourceCheck(false, model.ErrNoData.Error()), model.ErrNoData
}
