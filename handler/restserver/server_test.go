package restserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/n-creativesystem/rbns/handler/restserver"
	"github.com/n-creativesystem/rbns/handler/restserver/req"
	_ "github.com/n-creativesystem/rbns/infra"
	"github.com/n-creativesystem/rbns/infra/dao"
	_ "github.com/n-creativesystem/rbns/service"
	"github.com/n-creativesystem/rbns/tests"
	"github.com/stretchr/testify/assert"
)

func apiV1PostRequest(url string, obj interface{}) *http.Request {
	buf, _ := json.Marshal(obj)
	b := bytes.NewBuffer(buf)
	return httptest.NewRequest(http.MethodPost, path.Join("/api", "v1", url), b)
}

func TestServer(t *testing.T) {
	os.Setenv("MASTER_HOST", "api-rbac-postgres-dev")
	os.Setenv("MASTER_NAME", "postgres")
	os.Setenv("MASTER_USER", "rbac-user")
	os.Setenv("MASTER_PASSWORD", "rbac-user")
	os.Setenv("MASTER_PORT", "5432")
	dao.Register()
	_ = dao.New(
		dao.WithMigrationBack,
		dao.WithMigration,
		dao.WithDialector("postgres"),
		dao.WithMasterDSN(`host=${MASTER_HOST} user=${MASTER_USER} password=${MASTER_PASSWORD} dbname=${MASTER_NAME} port=${MASTER_PORT} sslmode=disable`),
		dao.WithSlaveDSN(`host=${MASTER_HOST} user=${MASTER_USER} password=${MASTER_PASSWORD} dbname=${MASTER_NAME} port=${MASTER_PORT} sslmode=disable`),
	)
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	srv := restserver.New(restserver.WithDebug, restserver.WithUI(false, "", "", false, "/"))
	cases := tests.Cases{
		{
			Name: "permissions",
			Fn: func(t *testing.T) {
				permissionCases := tests.Cases{
					{
						Name: "create",
						Fn: func(t *testing.T) {
							r = apiV1PostRequest("/permissions", &req.PermissionsCreateBody{
								Permissions: []req.PermissionCreateBody{
									{
										Name:        "create:test",
										Description: "create",
									},
								},
							})
							w = httptest.NewRecorder()
							srv.ServeHTTP(w, r)
							assert.Equal(t, http.StatusOK, w.Result().StatusCode)
						},
					},
					{
						Name: "find",
						Fn:   func(t *testing.T) {},
					},
				}
				permissionCases.Run(t)
			},
		},
	}

	cases.Run(t)
}
