package rdb_test

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
)

func testAddTenant(ctx context.Context) (*model.Tenant, error) {
	cmd := model.AddTenantCommand{
		Name: "test",
	}

	if err := bus.Dispatch(ctx, &cmd); err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
