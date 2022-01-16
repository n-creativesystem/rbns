package restserver

import (
	"net/http"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/ncsfw"
)

func (s *HTTPServer) AddAPIKey(c ncsfw.Context) error {
	cmd := model.AddApiKeyCommand{}
	if err := c.BindJSON(&cmd); err != nil {
		return response.ErrJson("bad request data", err)
	}
	if !cmd.Role.Valid() {
		return response.ErrJson("Invalid role specified", nil)
	}

	if !c.IsSignedIn() {
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "no login", nil)
	}

	user := c.GetLoginUser().(*model.LoginUser)
	tenant := c.GetTenant()
	cmd.ServiceAccountEmail = user.Email
	ctx := c.Request().Context()
	newKeyInfo, err := s.apiKeyService.Add(ctx, tenant, &cmd)
	if err != nil {
		return response.ErrJsonWithStatus(err.Code, err.Message, err.Error)
	}
	c.JSON(http.StatusOK, ncsfw.Map{
		"id":   cmd.Result.Id,
		"name": cmd.Result.Name,
		"key":  newKeyInfo.ClientSecret,
	})
	return nil
}

func (s *HTTPServer) DeleteAPIKey(c ncsfw.Context) error {
	strId := c.Param("id")
	if strId == "" {
		return response.ErrJson("bad request data", model.ErrRequired)
	}
	if err := s.apiKeyService.Delete(c.Request().Context(), strId); err != nil {
		return response.ErrJsonWithStatus(http.StatusInternalServerError, err.Error(), nil)
	}
	c.Status(http.StatusNoContent)
	return nil
}
