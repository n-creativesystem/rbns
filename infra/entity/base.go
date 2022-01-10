package entity

import (
	"context"
	"time"

	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"gorm.io/gorm"
)

type Model struct {
	ID        plugins.ID `gorm:"type:varchar(256);primarykey"`
	Tenant    string     `gorm:"type:varchar(256);primaryKey;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// func (m *Model) BeforeCreate(*gorm.DB) error {
// 	m.Generate()
// 	return nil
// }

func (m *Model) Generate() {
	m.ID.Generate()
}

func (m *Model) SetTenant(ctx context.Context) {
	if m.Tenant == "" {
		m.Tenant = contexts.FromTenantContext(ctx)
	}
}

// BeforeQuery Where条件に入る様に
func (m *Model) BeforeQuery(tx *gorm.DB) error {
	m.SetTenant(tx.Statement.Context)
	return nil
}
