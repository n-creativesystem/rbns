package restserver

import (
	"fmt"
	"net/http"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/handler/restserver/request"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
)

// addTenant テナントの新規作成
func (hs *HTTPServer) addTenant(c *contexts.Context) error {
	ctx := c.Request.Context()
	var reqTenant request.TenantPost
	err := c.BindJSON(&reqTenant)
	if err != nil {
		return err
	}

	cmd := model.AddTenantCommand{
		Name: reqTenant.Name,
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		hs.log.ErrorWithContext(ctx, err, "add tenant dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	publishCommand := model.AddTenantAndLoginUserCommand{
		Tenant:    cmd.Result,
		LoginUser: c.LoginUser,
	}

	if err := bus.PublishCtx(c, &publishCommand); err != nil {
		hs.log.ErrorWithContext(ctx, err, "add tenant publish error", "command", fmt.Sprintf("%+v", cmd))
		cmd := model.DeleteTenantCommand{
			ID: cmd.Result.ID,
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			hs.log.ErrorWithContext(ctx, err, "delete tenant dispatch error", "command", fmt.Sprintf("%+v", cmd))
			return err
		}
		return err
	}
	return nil
}

// getTenants ログインユーザーの保持しているテナント一覧の取得
func (hs *HTTPServer) getTenants(c *contexts.Context) error {
	ctx := c.Request.Context()
	query := model.GetLoginUserByIDQuery{
		ID: c.LoginUser.ID,
	}
	if err := bus.Dispatch(c, &query); err != nil {
		hs.log.ErrorWithContext(ctx, err, "get user dispatch error", "command", fmt.Sprintf("%+v", query))
		return err
	}
	tenants := response.NewTenants(query.Result.Tenants)
	body := response.NewResponse(tenants)
	c.JSON(http.StatusOK, &body)
	return nil
}
