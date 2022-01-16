package entity

import "fmt"

type ApiKey struct {
	Model
	Name string `gorm:"type:varchar(256);uniqueIndex"`
	Key  string `gorm:"column:key"`
	Role uint   `gorm:"role"`
}

func (ApiKey) IndexName(table, column string) string {
	return fmt.Sprintf("uq_%s_tanant", table)
}
