package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/handler/restserver/req"
	"github.com/n-creativesystem/rbns/handler/restserver/res"
	"github.com/n-creativesystem/rbns/service"
)

func init() {
	di.MustRegister(newUserHandler)
}

type userHandle interface {
	// Create is ユーザーを新規作成します。
	//
	// @body req.UserCreateBody
	create(*gin.Context)

	// findByKey is 組織IDとユーザーキーを条件にユーザーを取得します。
	//
	// @param id string
	//
	// @param key string
	findByKey(*gin.Context)

	// delete is 組織IDとユーザーキーを条件にユーザーを削除します。
	//
	// @param id string
	//
	// @param key string
	delete(*gin.Context)

	// addRole is 既存ユーザーに対してロールを追加します。
	//
	// @param id string
	//
	// @param key string
	//
	// @body req.UserUpdateBody
	addRole(*gin.Context)

	// deleteRole is 既存のユーザーからロールを削除します。
	//
	// @param id string
	//
	// @param key string
	//
	// @param roleId string
	deleteRole(*gin.Context)
}

type userHandler struct {
	svc service.UserService
}

func newUserHandler(svc service.UserService) userHandle {
	return &userHandler{svc: svc}
}

func (h *userHandler) create(c *gin.Context) {
	var req req.UserCreateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	roleIds := make([]string, len(req.Roles))
	for idx, role := range req.Roles {
		roleIds[idx] = role.Id
	}
	if err := h.svc.Create(c.Request.Context(), req.Key, req.OrganizationId, roleIds...); responseError(c, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *userHandler) findByKey(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	m, err := h.svc.FindByKey(c.Request.Context(), key, id)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewUser(*m))
}

func (h *userHandler) delete(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	if err := h.svc.Delete(c.Request.Context(), key, id); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *userHandler) addRole(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	var req req.UserUpdateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	roleIds := make([]string, len(req.Roles))
	for idx, role := range req.Roles {
		roleIds[idx] = role.Id
	}
	if err := h.svc.AddRole(c.Request.Context(), key, id, roleIds); responseError(c, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *userHandler) deleteRole(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	roleId := c.Param("roleId")
	if err := h.svc.DeleteRole(c.Request.Context(), key, id, []string{roleId}); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
