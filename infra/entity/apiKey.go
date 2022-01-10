package entity

import "fmt"

type ApiKey struct {
	Model
	Name string `gorm:"type:varchar(256);uniqueIndex"`
	Key  string `gorm:"column:key"`
	Role string `gorm:"role"`

	Tenants []Tenant `gorm:"many2many:tenant_api_keys;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ApiKey) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_tanant", table)
}
