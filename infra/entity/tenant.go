package entity

import (
	"fmt"
	"time"

	"github.com/n-creativesystem/rbns/infra/entity/plugins"
)

type Tenant struct {
	ID          plugins.ID `gorm:"type:varchar(256);primaryKey"`
	Name        string     `gorm:"not null;uniqueIndex"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Permissions 権限
	Permissions []Permission `gorm:"references:ID;foreignKey:Tenant;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// Roles ロール
	Roles []Role `gorm:"references:ID;foreignKey:Tenant;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// Organizations 組織
	Organizations []Organization `gorm:"references:ID;foreignKey:Tenant;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// Users ユーザー
	Users []User `gorm:"references:ID;foreignKey:Tenant;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// LoginUsers ログイン可能ユーザー
	LoginUsers []LoginUser `gorm:"many2many:tenant_login_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// ApiKeys  api key
	ApiKeys []ApiKey `gorm:"references:ID;foreignKey:Tenant;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Tenant) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_tanant", table)
}

type TenantUser struct {
	TenantID    plugins.ID `gorm:"type:varchar(256);primaryKey"`
	LoginUserID plugins.ID `gorm:"primaryKey"`

	// Users   *LoginUser `gorm:"foreignKey:ID;references:LoginUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// Tenants *Tenant    `gorm:"foreignKey:ID;references:TenantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
