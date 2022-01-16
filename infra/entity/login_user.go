package entity

import (
	"fmt"
	"time"
)

type LoginUser struct {
	OAuthId       string `gorm:"not null"`
	UserName      string `gorm:"not null"`
	Role          string `gorm:"not null"`
	OAuthName     string `gorm:"not null"`
	Password      string
	Email         string `gorm:"not null;primarykey"`
	OAuthToken    string
	SignupAllowed bool `gorm:"not null;default:false"`
	Tenant        string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Tenants []Tenant `gorm:"many2many:tenant_login_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (LoginUser) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_oauth_name", table)
}
