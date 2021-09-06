package service

import (
	"context"
	"errors"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/sirupsen/logrus"
)

type Resource interface {
	Find(ctx context.Context, id model.Key) (*model.Resource, error)
	FindAll(ctx context.Context) ([]*model.Resource, error)
	Exists(ctx context.Context, id model.Key) (bool, error)
	Delete(ctx context.Context, id model.Key) error
	Save(ctx context.Context, id model.Key, description string, permissions ...string) error
	Migration(ctx context.Context, id model.Key, description string, permissions ...string) error
	GetPermissions(ctx context.Context, id model.Key) (model.Permissions, error)
	AddPermissions(ctx context.Context, id model.Key, permissionIds []string) error
	DeletePermissions(ctx context.Context, id model.Key, permissionIds []string) error
}

type resourceService struct {
	reader repository.Reader
	writer repository.Writer
}

func NewResource(reader repository.Reader, writer repository.Writer) Resource {
	return &resourceService{
		reader: reader,
		writer: writer,
	}
}

func (svc *resourceService) Find(ctx context.Context, id model.Key) (*model.Resource, error) {
	strId := *id.Value()
	if _, err := model.NewKey(strId); err != nil {
		return nil, err
	}
	return svc.reader.Resource(ctx).Find(strId)
}

func (svc *resourceService) FindAll(ctx context.Context) ([]*model.Resource, error) {
	return svc.reader.Resource(ctx).FindAll()
}

func (svc *resourceService) Exists(ctx context.Context, id model.Key) (bool, error) {
	if v, err := svc.Find(ctx, id); err != nil && !errors.Is(err, model.ErrNoData) {
		return false, err
	} else {
		if v == nil {
			return false, nil
		}
		return true, nil
	}
}

func (svc *resourceService) Save(ctx context.Context, id model.Key, description string, permissions ...string) error {
	con := svc.reader
	permissionsRepo := con.Permission(ctx)
	mPermissions := make([]model.Permission, 0, len(permissions))
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
	return svc.writer.Do(ctx, func(tx repository.Transaction) error {
		return tx.Resource().Save(*id.Value(), description, mPermissions...)
	})
}

func (svc *resourceService) Delete(ctx context.Context, id model.Key) error {
	strId := *id.Value()
	return svc.writer.Do(ctx, func(tx repository.Transaction) error {
		return tx.Resource().Delete(strId)
	})
}

func (svc *resourceService) Migration(ctx context.Context, id model.Key, description string, permissions ...string) error {
	if description == "" {
		description = "resource auto migration create"
	}
	strId := *id.Value()
	// 既に存在するリソースの場合は何もしない
	if flg, _ := svc.Exists(ctx, id); flg {
		logrus.Debugf("exist id: %s", id)
		return nil
	}
	return svc.writer.Do(ctx, func(tx repository.Transaction) error {
		mPermissions := make(model.Permissions, 0, len(permissions))
		for _, permission := range permissions {
			pRepo := tx.Permission()
			name, err := model.NewName(permission)
			if err != nil {
				// 空文字は無視
				continue
			}
			if p, err := pRepo.FindByName(name); model.IsNoData(err) {
				p, err = pRepo.Create(name, "resource auto migration create")
				if err != nil {
					logrus.Errorf("permissions create err: %v", err)
					return err
				}
				mPermissions = append(mPermissions, *p)
			} else {
				mPermissions = append(mPermissions, *p)
			}
		}
		return tx.Resource().Save(strId, description, mPermissions...)
	})
}

func (svc *resourceService) GetPermissions(ctx context.Context, id model.Key) (model.Permissions, error) {
	strId := *id.Value()
	resource, err := svc.reader.Resource(ctx).Find(strId)
	if err != nil {
		return nil, err
	}
	permissions := resource.Permissions.Copy()
	return permissions, nil
}

func (svc *resourceService) AddPermissions(ctx context.Context, id model.Key, permissionIds []string) error {
	strId := *id.Value()
	return svc.writer.Do(ctx, func(tx repository.Transaction) error {
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
		return tx.Resource().AddPermission(strId, permissions)
	})
}

func (svc *resourceService) DeletePermissions(ctx context.Context, id model.Key, permissionIds []string) error {
	strId := *id.Value()
	if err := svc.writer.Do(ctx, func(tx repository.Transaction) error {
		for _, pId := range permissionIds {
			permissionId, err := model.NewID(pId)
			if err != nil {
				return err
			} else {
				if err := tx.Resource().DeletePermission(strId, permissionId); err != nil {
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
