package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/proto"
)

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
	reader repository.Reader
	writer repository.Writer
}

var _ RoleService = (*roleService)(nil)

func NewRoleService(reader repository.Reader, writer repository.Writer) RoleService {
	return &roleService{
		reader: reader,
		writer: writer,
	}
}

func (svc *roleService) Create(ctx context.Context, names, descriptions []string) (model.Roles, error) {
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
	err := svc.writer.Do(ctx, func(tx repository.Transaction) error {
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

func (svc *roleService) FindById(ctx context.Context, strId string) (*model.Role, error) {
	roleId, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	role, err := svc.reader.Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (svc *roleService) FindAll(ctx context.Context) (model.Roles, error) {
	roleRepo := svc.reader.Role(ctx)
	roles, err := roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (svc *roleService) Update(ctx context.Context, strId string, name, description string) error {
	mRole, err := model.NewRole(strId, name, description, nil)
	if err != nil {
		return err
	}
	if err := svc.writer.Do(ctx, func(tx repository.Transaction) error { return tx.Role().Update(mRole) }); err != nil {
		return err
	}
	return nil
}

func (svc *roleService) Delete(ctx context.Context, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	if err := svc.writer.Do(ctx, func(tx repository.Transaction) error { return tx.Role().Delete(id) }); err != nil {
		return err
	}
	return nil
}

func (svc *roleService) GetPermissions(ctx context.Context, strId string) (model.Permissions, error) {
	roleId, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	role, err := svc.reader.Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	permissions := role.GetPermissions().Copy()
	return permissions, nil
}

func (svc *roleService) AddPermissions(ctx context.Context, strId string, permissionIds []string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	err = svc.writer.Do(ctx, func(tx repository.Transaction) error {
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

func (svc *roleService) DeletePermissions(ctx context.Context, strId string, permissionIds []string) error {
	roleId, err := model.NewID(strId)
	if err != nil {
		return err
	}
	if err := svc.writer.Do(ctx, func(tx repository.Transaction) error {
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
