package service

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/proto"
)

func init() {
	di.MustRegister(newRoleService)
}

type RoleService interface {
	Create(ctx context.Context, names, descriptions []string) (model.Roles, error)
	FindById(ctx context.Context, strId string) (*model.Role, error)
	FindAll(ctx context.Context) (model.Roles, error)
	Update(ctx context.Context, strId string, name, description string) error
	Delete(ctx context.Context, strId string) error
	GetPermissions(ctx context.Context, strId string) (model.Permissions, error)
	AddPermissions(ctx context.Context, strId string, permissionIds []string) error
	DeletePermissions(ctx context.Context, strId string, permissionIds []string) error
}

type roleService struct {
	*proto.UnimplementedRoleServer
	repo repository.Repository
}

var _ RoleService = (*roleService)(nil)

func newRoleService(repo repository.Repository) RoleService {
	return &roleService{repo: repo}
}

func (srv *roleService) Create(ctx context.Context, names, descriptions []string) (model.Roles, error) {
	var out model.Roles
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
		roleRepo := tx.Role()
		roles, err := roleRepo.CreateBatch(mNames, mDescriptions)
		if err != nil {
			return err
		}
		out = make(model.Roles, len(roles))
		for idx, role := range roles {
			out[idx] = *role
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (srv *roleService) FindById(ctx context.Context, strId string) (*model.Role, error) {
	roleId, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	role, err := srv.repo.NewConnection().Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (srv *roleService) FindAll(ctx context.Context) (model.Roles, error) {
	roleRepo := srv.repo.NewConnection().Role(ctx)
	roles, err := roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (srv *roleService) Update(ctx context.Context, strId string, name, description string) error {
	mRole, err := model.NewRole(strId, name, description, nil)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Role().Update(mRole) }); err != nil {
		return err
	}
	return nil
}

func (srv *roleService) Delete(ctx context.Context, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Role().Delete(id) }); err != nil {
		return err
	}
	return nil
}

func (srv *roleService) GetPermissions(ctx context.Context, strId string) (model.Permissions, error) {
	roleId, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	role, err := srv.repo.NewConnection().Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	permissions := role.GetPermissions().Copy()
	return permissions, nil
}

func (srv *roleService) AddPermissions(ctx context.Context, strId string, permissionIds []string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		permissions := model.Permissions{}
		for _, permissionId := range permissionIds {
			if pId, err := model.NewID(permissionId); err == nil {
				if p, err := tx.Permission().FindByID(pId); err != nil {
					return err
				} else {
					permissions = append(permissions, *p)
				}
			} else {
				return err
			}
		}
		return tx.Role().AddPermission(id, permissions)
	})
	return err
}

func (srv *roleService) DeletePermissions(ctx context.Context, strId string, permissionIds []string) error {
	roleId, err := model.NewID(strId)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error {
		for _, id := range permissionIds {
			permissionId, err := model.NewID(id)
			if err != nil {
				return err
			} else {
				if err := tx.Role().DeletePermission(roleId, permissionId); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
