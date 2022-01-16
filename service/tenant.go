package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type Tenant interface {
	AddTenant(ctx context.Context, query *model.GetLoginUserByEmailQuery, cmd *model.AddTenantCommand) error
	GetTenant(ctx context.Context, query *model.GetLoginUserByEmailQuery) ([]model.Tenant, error)
}

type TenantImpl struct {
	log           logger.Logger
	telemetryFunc telemetryFunc
}

func NewTenantImpl() *TenantImpl {
	return &TenantImpl{
		log:           logger.New("tenant service"),
		telemetryFunc: createSpanWithPrefix("tenant service"),
	}
}

func (t *TenantImpl) AddTenant(ctx context.Context, query *model.GetLoginUserByEmailQuery, cmd *model.AddTenantCommand) (e error) {
	t.telemetryFunc(ctx, "add tenant", func(ctx context.Context) {
		if err := bus.Dispatch(ctx, &query); err != nil {
			t.log.ErrorWithContext(ctx, err, "query login user dispatch error", "query", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			t.log.ErrorWithContext(ctx, err, "add tenant dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (t *TenantImpl) GetTenant(ctx context.Context, query *model.GetLoginUserByEmailQuery) (out []model.Tenant, e error) {
	t.telemetryFunc(ctx, "get tenant", func(ctx context.Context) {
		if err := bus.Dispatch(ctx, &query); err != nil {
			t.log.ErrorWithContext(ctx, err, "get user dispatch error", "command", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out, e = query.Result.Tenants, nil
		return
	})
	return
}
