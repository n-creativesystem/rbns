package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
)

type RoleCheck func(ctx context.Context) error

var (
	Admin  = roleCheck(model.ROLE_ADMIN)
	Editor = roleCheck(model.ROLE_EDITOR)
	Viewer = roleCheck(model.ROLE_VIEWER)
)

func roleCheck(role model.RoleLevel) RoleCheck {
	return func(ctx context.Context) error {
		user, ok := FromCurrentUser(ctx)
		if !ok {
			return model.ErrForbidden
		}
		userRole, err := model.String2RoleLevel(user.Role)
		if err != nil {
			return model.ErrForbidden
		}
		if !role.IsLevelEnabled(userRole) {
			return model.ErrForbidden
		}
		return nil
	}
}
