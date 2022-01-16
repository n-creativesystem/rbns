package service

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	type test func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool
	cases := []struct {
		userRole model.RoleLevel
		admin    test
		editor   test
		viewer   test
	}{
		{
			userRole: model.ROLE_ADMIN,
			admin:    assert.NoError,
			editor:   assert.NoError,
			viewer:   assert.NoError,
		},
		{
			userRole: model.ROLE_EDITOR,
			admin:    assert.Error,
			editor:   assert.NoError,
			viewer:   assert.NoError,
		},
		{
			userRole: model.ROLE_VIEWER,
			admin:    assert.Error,
			editor:   assert.Error,
			viewer:   assert.NoError,
		},
	}
	for _, case_ := range cases {
		t.Run("user is "+case_.userRole.String(), func(t *testing.T) {
			user := &model.LoginUser{
				Email: "",
				Role:  case_.userRole.String(),
			}
			ctx := context.Background()
			ctx = SetCurrentUser(ctx, user)
			t.Run("admin", func(t *testing.T) {
				t.Helper()
				err := Admin(ctx)
				case_.admin(t, err)
			})
			t.Run("editor", func(t *testing.T) {
				t.Helper()
				err := Editor(ctx)
				case_.editor(t, err)
			})
			t.Run("viewer", func(t *testing.T) {
				t.Helper()
				err := Viewer(ctx)
				case_.viewer(t, err)
			})
		})
	}
}
