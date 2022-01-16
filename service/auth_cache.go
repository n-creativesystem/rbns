package service

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/n-creativesystem/rbns/cache"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthCache struct {
	Cfg   *config.Config
	Cache cache.Cache
}

var _ AuthService = (*AuthCache)(nil)

func NewAuthCache(cfg *config.Config, cache cache.Cache) *AuthCache {
	return &AuthCache{
		Cfg:   cfg,
		Cache: cache,
	}
}

func (a *AuthCache) SetToken(ctx context.Context, key string, user *model.LoginUser) error {
	ctx, span := tracer.Start(ctx, "service.SetToken", trace.WithAttributes(attribute.String("key", key)))
	defer span.End()
	buf, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return a.Cache.Set(key, buf)
}

func (a *AuthCache) LookupToken(ctx context.Context, key string) (model.LoginUser, error) {
	ctx, span := tracer.Start(ctx, "service.LookupToken", trace.WithAttributes(attribute.String("key", key)))
	defer span.End()
	var result model.LoginUser
	value, ok := a.Cache.Get(key)
	if !ok {
		return result, model.ErrNoDataFound
	}
	if err := json.NewDecoder(bytes.NewReader(value.([]byte))).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
