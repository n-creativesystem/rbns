package customs

import (
	"fmt"
	"strings"

	"gorm.io/gorm/schema"
)

type NamingStrategy struct {
	schema.NamingStrategy
}

var (
	_ schema.Namer = (*NamingStrategy)(nil)
)

func (ns NamingStrategy) IndexName(table, column string) string {
	column = strings.ToLower(column)
	if column == "tenant" || column == "name" {
		return fmt.Sprintf("idx_%s_tenant", table)
	}
	return ns.NamingStrategy.IndexName(table, column)
}
