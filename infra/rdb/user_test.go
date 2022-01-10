package rdb_test

import (
	"testing"
)

func TestUser(t *testing.T) {
	// tenant := "test"
	// ctx := context.Background()
	// ctx = contexts.ToTenantContext(ctx, tenant)
	// cases := tests.MocksByPostgres{
	// 	{
	// 		Name: "create",
	// 		Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
	// 			mock.ExpectBegin()
	// 			mock.ExpectExec(
	// 				regexp.QuoteMeta(`INSERT INTO "users" ("id","tenant","created_at","updated_at","name") VALUES ($1,$2,$3,$4,$5)`),
	// 			).WithArgs("", tests.IDs[0], tenant).WillReturnResult(sqlmock.NewResult(0, 1))
	// 			mock.ExpectCommit()
	// 			store := rdb.NewFactory(db, bus.GetBus())
	// 			return func(t *testing.T) {
	// 				orgId, _ := model.NewID(tests.IDs[0])
	// 				key, _ := model.NewKey(userKey)
	// 				cmd := model.AddUserCommand{}
	// 				err := store.AddUserCommand(ctx, &cmd)
	// 				assert.NoError(t, err)
	// 			}
	// 		},
	// 	},
	// 	{
	// 		Name: "add role",
	// 		Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
	// 			const (
	// 				userKey = "user1"
	// 			)
	// 			mock.ExpectBegin()
	// 			mock.ExpectExec(
	// 				regexp.QuoteMeta(`INSERT INTO "user_roles" ("user_key","organization_id","role_id","organization_tenant","role_tenant") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10) ON CONFLICT DO NOTHING`),
	// 			).WithArgs(userKey, tests.IDs[0], tests.IDs[1], tenant, tenant, userKey, tests.IDs[0], tests.IDs[2], tenant, tenant).WillReturnResult(sqlmock.NewResult(0, 2))
	// 			mock.ExpectCommit()
	// 			store := rdb.NewFactory(db, bus.GetBus())
	// 			return func(t *testing.T) {
	// 				orgId, _ := model.NewID(tests.IDs[0])
	// 				userKey, _ := model.NewKey(userKey)
	// 				role, _ := model.NewID(tests.IDs[1])
	// 				role2, _ := model.NewID(tests.IDs[2])
	// 				cmd := model.AddUserRoleCommand{
	// 					OrgID: orgId,
	// 					Key:   userKey,
	// 					RoleIDs: []model.ID{role, role2},
	// 				}
	// 				err := store.AddUserRoleCommand(ctx, &cmd)
	// 				assert.NoError(t, err)
	// 			}
	// 		},
	// 	},
	// 	{
	// 		Name: "delete role",
	// 		Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
	// 			mock.ExpectBegin()
	// 			mock.ExpectExec(
	// 				regexp.QuoteMeta(`DELETE FROM "user_roles" WHERE "user_roles"."user_key" = $1 AND "user_roles"."organization_id" = $2 AND "user_roles"."role_id" = $3 AND "user_roles"."organization_tenant" = $4`),
	// 			).WithArgs("user1", tests.IDs[0], tests.IDs[1], tenant).WillReturnResult(sqlmock.NewResult(0, 1))
	// 			mock.ExpectCommit()
	// 			store := rdb.NewFactory(db, bus.GetBus())
	// 			return func(t *testing.T) {
	// 				orgId, _ := model.NewID(tests.IDs[0])
	// 				userKey, _ := model.NewKey("user1")
	// 				role, _ := model.NewID(tests.IDs[1])
	// 				cmd := model.DeleteUserRoleCommand{
	// 					OrgID: orgId,
	// 					Key:   userKey,
	// 					RoleIDs: []model.ID{role},
	// 				}
	// 				err := store.DeleteUserRoleCommand(ctx, &cmd)
	// 				assert.NoError(t, err)
	// 			}
	// 		},
	// 	},
	// 	{
	// 		Name: "findByKey",
	// 		Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
	// 			usersRow := sqlmock.NewRows([]string{"organization_id", "key"}).AddRow(tests.IDs[0], "user1")
	// 			orgRow := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[0], "test", "organization test")
	// 			userRoleRow := sqlmock.NewRows([]string{"user_key", "organization_id", "role_id"}).AddRow("user1", tests.IDs[0], tests.IDs[1]).AddRow("user1", tests.IDs[0], tests.IDs[2])
	// 			roleRow := sqlmock.
	// 				NewRows([]string{"id", "name", "description"}).
	// 				AddRow(tests.IDs[1], "admin", "administrator").
	// 				AddRow(tests.IDs[2], "view", "viewer")
	// 			rolePermissionRow := sqlmock.NewRows([]string{"role_id", "permission_id"}).AddRow(tests.IDs[1], tests.IDs[3]).AddRow(tests.IDs[2], tests.IDs[4])
	// 			permissionRow := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[3], "create:permission", "create").AddRow(tests.IDs[4], "read:permission", "read")
	// 			// users find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."key" = $1 AND "users"."organization_id" = $2 AND "users"."organization_tenant" = $3`),
	// 			).WithArgs("user1", tests.IDs[0], tenant).WillReturnRows(usersRow)

	// 			// organizations find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE ("organizations"."id","organizations"."tenant") IN (($1,$2)) AND "organizations"."tenant" = $3`),
	// 			).WithArgs(tests.IDs[0], tenant, tenant).WillReturnRows(orgRow)

	// 			// user_roles find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE ("user_roles"."user_key","user_roles"."organization_id","user_roles"."organization_tenant") IN (($1,$2,$3)) AND "user_roles"."organization_tenant" = $4`),
	// 			).WithArgs("user1", tests.IDs[0], tenant, tenant).WillReturnRows(userRoleRow)

	// 			// roles find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "roles" WHERE ("roles"."id","roles"."tenant") IN (($1,$2),($3,$4)) AND "roles"."tenant" = $5`),
	// 			).WithArgs(tests.IDs[1], tenant, tests.IDs[2], tenant, tenant).WillReturnRows(roleRow)

	// 			// role_permissions find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE ("role_permissions"."role_id","role_permissions"."role_tenant") IN (($1,$2),($3,$4)) AND "role_permissions"."role_tenant" = $5`),
	// 			).WithArgs(tests.IDs[1], tenant, tests.IDs[2], tenant, tenant).WillReturnRows(rolePermissionRow)

	// 			// permissions find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE ("permissions"."id","permissions"."tenant") IN (($1,$2),($3,$4)) AND "permissions"."tenant" = $5`),
	// 			).WithArgs(tests.IDs[3], tenant, tests.IDs[4], tenant, tenant).WillReturnRows(permissionRow)

	// 			store := rdb.NewFactory(db, bus.GetBus())
	// 			return func(t *testing.T) {
	// 				orgId, _ := model.NewID(tests.IDs[0])
	// 				userKey, _ := model.NewKey("user1")
	// 				query := model.GetUserByKeyQuery{
	// 					OrgID: orgId,
	// 					Key:   userKey,
	// 				}
	// 				err := store.GetUserByKeyQuery(ctx, &query)
	// 				assert.NoError(t, err)
	// 				u := query.Result
	// 				assert.Equal(t, userKey.Value().String(), u.GetKey())
	// 				assert.NotEmpty(t, u.GetRole())
	// 				assert.NotEmpty(t, u.GetPermission())
	// 				assert.Equal(t, tests.IDs[1], u.GetRole()[0].GetID().String())
	// 				assert.Equal(t, tests.IDs[2], u.GetRole()[1].GetID().String())
	// 				assert.Equal(t, tests.IDs[3], u.GetPermission()[0].GetID().String())
	// 				assert.Equal(t, tests.IDs[4], u.GetPermission()[1].GetID().String())
	// 			}
	// 		},
	// 	},
	// 	{
	// 		Name: "findAll",
	// 		Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
	// 			user1 := "user1"
	// 			user2 := "user2"
	// 			orgId := tests.IDs[0]
	// 			admin := tests.IDs[1]
	// 			view := tests.IDs[2]
	// 			pCreate := tests.IDs[3]
	// 			pRead := tests.IDs[4]
	// 			usersRow := sqlmock.NewRows([]string{"organization_id", "key"}).
	// 				AddRow(orgId, user1).
	// 				AddRow(orgId, user2)
	// 			orgRow := sqlmock.NewRows([]string{"id", "name", "description"}).
	// 				AddRow(orgId, "test", "organization test")
	// 			userRoleRow := sqlmock.NewRows([]string{"user_key", "organization_id", "role_id"}).
	// 				AddRow(user1, orgId, admin).AddRow(user1, orgId, view).
	// 				AddRow("user2", orgId, view)
	// 			roleRow := sqlmock.
	// 				NewRows([]string{"id", "name", "description"}).
	// 				AddRow(admin, "admin", "administrator").
	// 				AddRow(view, "view", "viewer")
	// 			rolePermissionRow := sqlmock.NewRows([]string{"role_id", "permission_id"}).
	// 				AddRow(admin, pCreate).
	// 				AddRow(view, pRead)
	// 			permissionRow := sqlmock.NewRows([]string{"id", "name", "description"}).
	// 				AddRow(pCreate, "create:permission", "create").
	// 				AddRow(pRead, "read:permission", "read")
	// 			// users find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."organization_id" = $1 AND "users"."organization_tenant" = $2`),
	// 			).WithArgs(orgId, tenant).WillReturnRows(usersRow)

	// 			// organizations find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE ("organizations"."id","organizations"."tenant") IN (($1,$2)) AND "organizations"."tenant" = $3`),
	// 			).WithArgs(orgId, tenant, tenant).WillReturnRows(orgRow)

	// 			// user_roles find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE ("user_roles"."user_key","user_roles"."organization_id","user_roles"."organization_tenant") IN (($1,$2,$3),($4,$5,$6)) AND "user_roles"."organization_tenant" = $7`),
	// 			).WithArgs(user1, orgId, tenant, user2, orgId, tenant, tenant).WillReturnRows(userRoleRow)

	// 			// roles find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "roles" WHERE ("roles"."id","roles"."tenant") IN (($1,$2),($3,$4)) AND "roles"."tenant" = $5`),
	// 			).WithArgs(admin, tenant, view, tenant, tenant).WillReturnRows(roleRow)

	// 			// role_permissions find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE ("role_permissions"."role_id","role_permissions"."role_tenant") IN (($1,$2),($3,$4)) AND "role_permissions"."role_tenant" = $5`),
	// 			).WithArgs(admin, tenant, view, tenant, tenant).WillReturnRows(rolePermissionRow)

	// 			// permissions find
	// 			mock.ExpectQuery(
	// 				regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE ("permissions"."id","permissions"."tenant") IN (($1,$2),($3,$4)) AND "permissions"."tenant" = $5`),
	// 			).WithArgs(pCreate, tenant, pRead, tenant, tenant).WillReturnRows(permissionRow)

	// 			store := rdb.NewFactory(db, bus.GetBus())
	// 			return func(t *testing.T) {
	// 				orgId, _ := model.NewID(tests.IDs[0])
	// 				query := model.GetUserQuery{
	// 					OrgId: orgId,
	// 				}
	// 				err := store.GetUserQuery(ctx, &query)
	// 				if !assert.NoError(t, err) {
	// 					return
	// 				}
	// 				us := query.Result
	// 				u := us[0]
	// 				assert.Equal(t, user1, u.GetKey())
	// 				assert.NotEmpty(t, u.GetRole())
	// 				assert.NotEmpty(t, u.GetPermission())
	// 				assert.Equal(t, 2, len(u.GetRole()))
	// 				assert.Equal(t, 2, len(u.GetPermission()))
	// 				assert.Equal(t, admin, u.GetRole()[0].GetID().String())
	// 				assert.Equal(t, view, u.GetRole()[1].GetID().String())
	// 				assert.Equal(t, pCreate, u.GetPermission()[0].GetID().String())
	// 				assert.Equal(t, pRead, u.GetPermission()[1].GetID().String())
	// 				u = us[1]
	// 				assert.Equal(t, user2, u.GetKey())
	// 				assert.NotEmpty(t, u.GetRole())
	// 				assert.NotEmpty(t, u.GetPermission())
	// 				assert.Equal(t, 1, len(u.GetRole()))
	// 				assert.Equal(t, 1, len(u.GetPermission()))
	// 				assert.Equal(t, view, u.GetRole()[0].GetID().String())
	// 				assert.Equal(t, pRead, u.GetPermission()[0].GetID().String())
	// 			}
	// 		},
	// 	},
	// }
	// cases.Run(t)
}
