package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
)

type currentUser struct{}

func SetCurrentUser(ctx context.Context, user *model.LoginUser) context.Context {
	return context.WithValue(ctx, currentUser{}, user)
}

func FromCurrentUser(ctx context.Context) (*model.LoginUser, bool) {
	v, ok := ctx.Value(currentUser{}).(*model.LoginUser)
	return v, ok
}

type AuthService interface {
	SetToken(ctx context.Context, key string, user *model.LoginUser) error
	LookupToken(ctx context.Context, key string) (model.LoginUser, error)
	// LookupApiKey(key string)
}
