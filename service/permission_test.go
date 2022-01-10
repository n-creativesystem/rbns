package service

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/mock"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"github.com/n-creativesystem/rbns/tracer"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPermission(t *testing.T) {
	const tenant = "test"
	ctx := context.Background()
	ctx = contexts.ToTenantContext(ctx, tenant)
	trace_, _ := tracer.InitOpenTelemetry("github.com/n-creativesystem/rbns/test/service")
	defer trace_.Cleanup(ctx)
	ctx, span := tracer.Start(ctx, "permission service")
	defer span.End()
	mock.NewCase(mock.PostgreSQL, "permission_svc").Set(mock.Case{
		Name: "permission service test",
		Fn: func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T) {
			svc := NewPermissionService()
			expecteds := []struct {
				name, description string
			}{
				{
					name:        "create:user",
					description: "create user permission",
				},
				{
					name:        "update:user",
					description: "update user permisssion",
				},
				{
					name:        "read:user",
					description: "read user permission",
				},
				{
					name:        "delete:user",
					description: "delete user permission",
				},
			}
			return func(t *testing.T) {
				var created []model.Permission
				t.Run("create", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "create")
					defer span.End()
					names := []string{"create:user", "update:user", "read:user", "delete:user"}
					descriptions := []string{"create user permission", "update user permisssion", "read user permission", "delete user permission"}
					out, err := svc.Create(ctx, names, descriptions)
					assert.NoError(t, err)
					for idx, entity := range out {
						expected := expecteds[idx]
						assert.NotEmpty(t, entity.ID.String())
						assert.Equal(t, expected.name, entity.Name)
						assert.Equal(t, expected.description, entity.Description)
					}
					created = out
				})
				t.Run("create duplicate", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "create duplicate")
					defer span.End()
					names := []string{"create:user"}
					descriptions := []string{"create user permission"}
					out, err := svc.Create(ctx, names, descriptions)
					assert.Error(t, err)
					assert.ErrorIs(t, err, model.ErrAlreadyExists)
					assert.Nil(t, out)
				})
				t.Run("find by id", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "find by id")
					defer span.End()
					out, err := svc.FindById(ctx, created[0].ID.String())
					assert.NoError(t, err)
					assert.Equal(t, created[0].ID.String(), out.ID.String())
					assert.Equal(t, created[0].Name, out.Name)
					assert.Equal(t, created[0].Description, out.Description)
				})
				t.Run("find all", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "find all")
					defer span.End()
					out, err := svc.FindAll(ctx)
					assert.NoError(t, err)
					for idx, entity := range out {
						expected := created[idx]
						assert.Equal(t, expected.ID.String(), entity.ID.String())
						assert.Equal(t, expected.Name, entity.Name)
						assert.Equal(t, expected.Description, entity.Description)
					}
				})
				t.Run("update", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "update")
					defer span.End()
					out, err := svc.Update(ctx, created[0].ID.String(), "read:permission", "read user permission")
					assert.NoError(t, err)
					assert.Equal(t, "read:permission", out.Name)
					assert.Equal(t, "read user permission", out.Description)
				})
				t.Run("delete", func(t *testing.T) {
					ctx, span := tracer.Start(ctx, "delete")
					defer span.End()
					err := svc.Delete(ctx, created[0].ID.String())
					assert.NoError(t, err)
					out, err := svc.FindById(ctx, created[0].ID.String())
					assert.Error(t, err)
					assert.ErrorIs(t, err, model.ErrNoData)
					assert.Nil(t, out)
				})
			}
		},
	}).Run(t)
}
