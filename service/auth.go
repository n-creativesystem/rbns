package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
)

type AuthService interface {
	SetToken(ctx context.Context, key string, user *model.LoginUser) error
	LookupToken(ctx context.Context, key string) (model.LoginUser, error)
	// LookupApiKey(key string)
}
