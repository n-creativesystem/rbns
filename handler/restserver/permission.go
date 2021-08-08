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
	di.MustRegister(newPermissionHandler)
}

type permissionHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
}

type permissionHandler struct {
	svc service.PermissionService
}

func newPermissionHandler(svc service.PermissionService) permissionHandle {
	return &permissionHandler{svc: svc}
}

func (h *permissionHandler) create(c *gin.Context) {
	var req req.PermissionsCreateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	names := make([]string, len(req.Permissions))
	descriptions := make([]string, len(req.Permissions))
	for idx, p := range req.Permissions {
		names[idx] = p.Name
		descriptions[idx] = p.Description
	}
	m, err := h.svc.Create(c.Request.Context(), names, descriptions)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewPermissions(m))
}

func (h *permissionHandler) findById(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.FindById(c.Request.Context(), id)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewPermission(*m))
}

func (h *permissionHandler) findAll(c *gin.Context) {
	m, err := h.svc.FindAll(c.Request.Context())
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewPermissions(m))
}

func (h *permissionHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req req.PermissionUpdateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, req.Name, req.Description); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *permissionHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
