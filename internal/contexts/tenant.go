package contexts

import (
	"context"
	"errors"

	"github.com/n-creativesystem/rbns/logger"
)

type tenantKey struct{}

var key = tenantKey{}

func ToTenantContext(ctx context.Context, tenant string) context.Context {
	if tenant == "" {
		err := errors.New("Tenant is empty")
		logger.PanicWithContext(ctx, err, "ToTenantContext")
	}
	return context.WithValue(ctx, key, tenant)
}

func FromTenantContext(ctx context.Context) string {
	tenant, _ := ctx.Value(key).(string)
	return tenant
}
