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
	di.MustRegister(newRoleHandler)
}

type roleHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
	getPermissions(*gin.Context)
	addPermissions(*gin.Context)
	deletePermissions(*gin.Context)
}

type roleHandler struct {
	svc service.RoleService
}

func newRoleHandler(svc service.RoleService) roleHandle {
	return &roleHandler{svc: svc}
}

func (h *roleHandler) create(c *gin.Context) {
	var req req.RolesCreateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	names := make([]string, len(req.Roles))
	descriptions := make([]string, len(req.Roles))
	for idx, role := range req.Roles {
		names[idx] = role.Name
		descriptions[idx] = role.Description
	}
	m, err := h.svc.Create(c.Request.Context(), names, descriptions)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewRoles(m))
}

func (h *roleHandler) findById(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.FindById(c.Request.Context(), id)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewRole(*m))
}

func (h *roleHandler) findAll(c *gin.Context) {
	m, err := h.svc.FindAll(c.Request.Context())
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewRoles(m))
}

func (h *roleHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req req.RoleUpdateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, req.Name, req.Description); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *roleHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *roleHandler) getPermissions(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.GetPermissions(c.Request.Context(), id)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewPermissions(m))
}

func (h *roleHandler) addPermissions(c *gin.Context) {
	id := c.Param("id")
	var req req.RolePermission
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	permissionIds := make([]string, len(req.Permissions))
	for idx, p := range req.Permissions {
		permissionIds[idx] = p.Id
	}
	err := h.svc.AddPermissions(c.Request.Context(), id, permissionIds)
	if responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *roleHandler) deletePermissions(c *gin.Context) {
	roleId := c.Param("id")
	permissionId := c.Param("permissionId")
	if err := h.svc.DeletePermissions(c.Request.Context(), roleId, []string{permissionId}); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
