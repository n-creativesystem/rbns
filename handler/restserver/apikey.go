package restserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/internal/apikeygen"
)

func (s *HTTPServer) AddAPIKey(c *contexts.Context) error {
	cmd := model.AddApiKeyCommand{}
	if err := c.BindJSON(&cmd); err != nil {
		return response.ErrJson("bad request data", err)
	}
	if !cmd.Role.Valid() {
		return response.ErrJson("Invalid role specified", nil)
	}

	if !c.IsSignedIn {
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "no login", nil)
	}
	user := c.LoginUser
	tenant := c.Tenant
	cmd.ServiceAccountID = user.ID
	newKeyInfo, err := apikeygen.New(cmd.ServiceAccountID, tenant, cmd.Name)
	if err != nil {
		return response.ErrJsonWithStatus(http.StatusInternalServerError, "Generating API key failed", err)
	}
	cmd.HashedKey = newKeyInfo.HashedKey
	if err := bus.Dispatch(c.Request.Context(), &cmd); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return response.ErrJsonWithStatus(http.StatusConflict, err.Error(), nil)
		}
		return response.ErrJsonWithStatus(http.StatusInternalServerError, "Failed to add API Key", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   cmd.Result.Id,
		"name": cmd.Result.Name,
		"key":  newKeyInfo.ClientSecret,
	})
	return nil
}

func (s *HTTPServer) DeleteAPIKey(c *contexts.Context) error {
	strId := c.Param("id")
	if strId == "" {
		return response.ErrJson("bad request data", model.ErrRequired)
	}
	id, _ := model.NewID(strId)
	cmd := model.DeleteAPIKeyCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
	}
	if err := bus.Dispatch(c, &cmd); err != nil {
		return response.ErrJsonWithStatus(http.StatusInternalServerError, "Failed to delete API key", err)
	}
	c.Status(http.StatusNoContent)
	return nil
}
