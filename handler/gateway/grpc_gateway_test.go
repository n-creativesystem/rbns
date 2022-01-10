package gateway

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/n-creativesystem/rbns/logger"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	var mux http.Handler
	var err error
	if mux, err = New(nil); err != nil {
		assert.NoError(t, err)
	}
	urls := []struct {
		method string
		url    string
		value  string
		status int
		result func(w *httptest.ResponseRecorder, req *http.Request)
	}{
		{
			method: http.MethodGet,
			url:    "/api/v1/swagger.json",
			status: http.StatusOK,
			result: func(w *httptest.ResponseRecorder, req *http.Request) {
				logger.InfoWithContext(req.Context(), w.Body.String())
			},
		},
		{
			method: http.MethodDelete,
			status: http.StatusOK,
			url:    "/api/v1/organizations/organizationId/users/userKey/roles/roleId",
		},
		{
			method: http.MethodPost,
			url:    "/api/v1/roles",
			value:  `{"roles":[{"name":"test","description":"テスト"}]}`,
			status: http.StatusOK,
		},
	}
	for _, sts := range urls {
		var req *http.Request
		var w *httptest.ResponseRecorder
		var reader io.Reader
		if sts.value != "" {
			reader = strings.NewReader(sts.value)
		}
		req = httptest.NewRequest(sts.method, sts.url, reader)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, sts.status, w.Result().StatusCode)
		if sts.result != nil {
			sts.result(w, req)
		}
	}
}
