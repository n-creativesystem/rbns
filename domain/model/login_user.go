package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/n-creativesystem/rbns/internal/utils"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/utilsconv"
	"golang.org/x/oauth2"
)

type LoginUser struct {
	// ID rbns id
	// ID string
	// OAuthID oauth provider id
	OAuthID string
	// UserName oauth username
	UserName string
	// Email oauth email
	Email string
	// Role user role
	Role   string
	Groups []string
	// Tenant 現在処理しているテナントが設定される
	tenant string

	OAuthToke *oauth2.Token

	password string `json:"-"`
	// oauthName oauth provider name
	oauthName string `json:"-"`

	err error `json:"-"`

	Tenants []Tenant `json:"-"`
}

func (user *LoginUser) SetOAuthName(name string) *LoginUser {
	user.oauthName = fmt.Sprintf("oauth_%s", name)
	return user
}

func (user *LoginUser) GetOAuthName() string {
	return strings.TrimPrefix(user.oauthName, "oauth_")
}

func (user *LoginUser) SetPassword(password string) *LoginUser {
	if password == "" && user.oauthName == "local" {
		user.error(errors.New("Password required"))
		return user
	}
	if password != "" {
		v, _ := utils.EncodePassword(password, user.Email)
		if v != password {
			// 暗号化されていない場合は暗号化した内容でセット
			password = v
		}
		user.password = password
	}
	return user
}

func (user *LoginUser) GetPassword() string {
	return user.password
}

func (user *LoginUser) GetToken() string {
	if user.OAuthToke != nil {
		buf, _ := json.Marshal(user.OAuthToke)
		return utilsconv.BytesToBase64(buf)
	}
	return ""
}

func (user *LoginUser) SetOAuthToken(token string) *LoginUser {
	var oauthToken oauth2.Token
	if v, err := utilsconv.Base64ToByte(token); err == nil {
		if !user.error(json.Unmarshal(v, &oauthToken)) {
			user.OAuthToke = &oauthToken
		}
	} else {
		logger.Error(err, "LoginUser.SetOauthToken")
	}
	return user
}

func (user *LoginUser) error(err error) bool {
	if err != nil {
		user.err = err
		return true
	}
	return false
}

func (user *LoginUser) Valid() bool {
	if user.err != nil {
		return false
	}
	// if !user.OAuthToke.Valid() {
	// 	return false
	// }
	_, err := String2RoleLevel(user.Role)
	if err != nil {
		return false
	}
	return true
}

func (user *LoginUser) GetTenant() string {
	if user.tenant == "" {
		if len(user.Tenants) > 0 {
			user.tenant = user.Tenants[0].ID
		}
	}
	return user.tenant
}

func (user *LoginUser) SetTenant(tenant string) *LoginUser {
	user.tenant = tenant
	return user
}

// IsVerify テナント登録が完了しているかどうか
func (user *LoginUser) IsVerify() bool {
	return user.GetTenant() != ""
}

func (user *LoginUser) Error() string {
	if user.err != nil {
		return user.err.Error()
	}
	return ""
}

func (user *LoginUser) Serialize() string {
	buf, _ := json.Marshal(user)
	return utilsconv.BytesToString(buf)
}

func (user *LoginUser) Deserialize(jsonString string) error {
	return json.Unmarshal(utilsconv.StringToBytes(jsonString), user)
}

type UpsertLoginUserCommand struct {
	User          *LoginUser
	SignupAllowed bool
}

type AddTenantLoginUserCommand struct {
	User   *LoginUser
	Tenant *Tenant
}

type GetLoginUserByEmailQuery struct {
	Email string

	Result *LoginUser
}

func (query *GetLoginUserByEmailQuery) Valid() error {
	if query.Email == "" {
		return ErrRequired
	}
	return nil
}

type GetLoginUserQuery struct {
	Result []*LoginUser
}
