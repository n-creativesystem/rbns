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
	di.MustRegister(newOrganizationHandler)
}

type organizationHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
}

type organizationHandler struct {
	svc service.OrganizationService
}

func newOrganizationHandler(svc service.OrganizationService) organizationHandle {
	return &organizationHandler{svc: svc}
}

func (h *organizationHandler) create(c *gin.Context) {
	var req req.OrganizationCreateBody
	if err := c.Bind(&req); requestError(c, err, body) {
		return
	}
	m, err := h.svc.Create(c.Request.Context(), req.Name, req.Description)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewOrganization(*m))
}

func (h *organizationHandler) findById(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.FindById(c.Request.Context(), id)
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewOrganization(*m))
}

func (h *organizationHandler) findAll(c *gin.Context) {
	m, err := h.svc.FindAll(c.Request.Context())
	if responseError(c, err) {
		return
	}
	c.JSON(http.StatusOK, res.NewOrganizations(m))
}

func (h *organizationHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req req.OrganizationUpdateBody
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, req.Name, req.Description); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *organizationHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
