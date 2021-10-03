package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
)

type UserService interface {
	Create(ctx context.Context, userKey, organizationId string, roleIds ...string) error
	Delete(ctx context.Context, userKey, organizationId string) error
	FindByKey(ctx context.Context, userKey, organizationId string) (*model.User, error)
	AddRole(ctx context.Context, userKey, organizationId string, roleIds []string) error
	DeleteRole(ctx context.Context, userKey, organizationId string, roleIds []string) error
}

type userService struct {
	reader repository.Reader
	writer repository.Writer
}

var _ UserService = (*userService)(nil)

func NewUserService(reader repository.Reader, writer repository.Writer) UserService {
	return &userService{
		reader: reader,
		writer: writer,
	}
}

// User
func (svc *userService) Create(ctx context.Context, userKey, organizationId string, roleIds ...string) error {
	orgId, err := model.NewID(organizationId)
	if err != nil {
		return err
	}
	user, err := model.NewUser(userKey, nil, nil)
	if err != nil {
		return err
	}
	err = svc.writer.Do(ctx, func(tx repository.Transaction) error {
		mRole := model.Roles{}
		for _, roleId := range roleIds {
			if id, err := model.NewID(roleId); err == nil {
				if r, err := tx.Role().FindByID(id); err == nil {
					mRole = append(mRole, *r)
				} else {
					return err
				}
			} else {
				return err
			}
		}
		user.AddRole(mRole...)
		userRepo := tx.User()
		_, err := userRepo.Create(orgId, user)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) Delete(ctx context.Context, userKey, organizationId string) error {
	orgId, err := model.NewID(organizationId)
	if err != nil {
		return err
	}
	key, err := model.NewKey(userKey)
	if err != nil {
		return err
	}
	err = svc.writer.Do(ctx, func(tx repository.Transaction) error {
		return tx.User().Delete(orgId, key)
	})
	return err
}

func (svc *userService) FindByKey(ctx context.Context, userKey, organizationId string) (*model.User, error) {
	orgId, err := model.NewID(organizationId)
	if err != nil {
		return nil, err
	}
	key, err := model.NewKey(userKey)
	if err != nil {
		return nil, err
	}
	userRepo := svc.reader.User(ctx)
	u, err := userRepo.FindByKey(orgId, key)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (svc *userService) AddRole(ctx context.Context, userKey, organizationId string, roleIds []string) error {
	orgId, err := model.NewID(organizationId)
	if err != nil {
		return err
	}
	key, err := model.NewKey(userKey)
	if err != nil {
		return err
	}
	err = svc.writer.Do(ctx, func(tx repository.Transaction) error {
		mRole := model.Roles{}
		for _, roleId := range roleIds {
			if id, err := model.NewID(roleId); err == nil {
				if r, err := tx.Role().FindByID(id); err == nil {
					mRole = append(mRole, *r)
				} else {
					return err
				}
			} else {
				return err
			}
		}
		return tx.User().AddRole(orgId, key, mRole...)
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) DeleteRole(ctx context.Context, userKey, organizationId string, roleIds []string) error {
	orgId, err := model.NewID(organizationId)
	if err != nil {
		return err
	}
	key, err := model.NewKey(userKey)
	if err != nil {
		return err
	}
	err = svc.writer.Do(ctx, func(tx repository.Transaction) error {
		for _, roleId := range roleIds {
			if id, err := model.NewID(roleId); err != nil {
				return err
			} else {
				if err := tx.User().DeleteRole(orgId, key, id); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
