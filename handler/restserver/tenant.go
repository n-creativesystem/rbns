package restserver

import (
	"net/http"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/request"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/ncsfw"
)

// addTenant テナントの新規作成
func (hs *HTTPServer) addTenant(c ncsfw.Context) error {
	ctx := c.Request().Context()
	var reqTenant request.TenantPost
	err := c.BindJSON(&reqTenant)
	if err != nil {
		return err
	}

	cmd := model.AddTenantCommand{
		Name:        model.TenantName(reqTenant.Name),
		Description: reqTenant.Description,
	}
	query := model.GetLoginUserByEmailQuery{
		Email: c.GetLoginUser().(*model.LoginUser).Email,
	}
	if err := hs.tenantService.AddTenant(ctx, &query, &cmd); err != nil {
		return err
	}
	c.Status(http.StatusCreated)
	return nil
}

// getTenants ログインユーザーの保持しているテナント一覧の取得
func (hs *HTTPServer) getTenants(c ncsfw.Context) error {
	ctx := c.Request().Context()
	query := model.GetLoginUserByEmailQuery{
		Email: c.GetLoginUser().(*model.LoginUser).Email,
	}
	results, err := hs.tenantService.GetTenant(ctx, &query)
	if err != nil {
		return err
	}
	tenants := response.NewTenants(results)
	body := response.NewResponse(tenants)
	c.JSON(http.StatusOK, &body)
	return nil
}
