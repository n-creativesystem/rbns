package restserver

import (
	"net/http"

	"github.com/n-creativesystem/rbns/domain/dtos"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw"
)

func (s *HTTPServer) Index(c ncsfw.Context) error {
	var loginUser model.LoginUser
	tenants := make([]dtos.Tenant, 0, 10)
	if c.IsSignedIn() {
		loginUser = *c.GetLoginUser().(*model.LoginUser)
		for _, t := range loginUser.Tenants {
			tenants = append(tenants, dtos.Tenant{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
			})
		}
	}
	if !loginUser.Valid() {
		loginUser = model.LoginUser{}
		c.SetLoginUser(&loginUser)
	}
	currentUser := dtos.CurrentUser{
		ID:         loginUser.OAuthID,
		UseName:    loginUser.UserName,
		Email:      loginUser.Email,
		Role:       loginUser.Role,
		Groups:     loginUser.Groups,
		IsSignedIn: c.IsSignedIn(),
		IsTenant:   loginUser.IsVerify(),
		Tenant:     loginUser.GetTenant(),
	}
	c.HTML(http.StatusOK, "index.html", ncsfw.Map{
		"BaseURL": s.Cfg.RootURL.String(),
		"SubPath": s.Cfg.SubPath,
		"User":    currentUser,
		"Tenants": tenants,
	})
	return nil
}
