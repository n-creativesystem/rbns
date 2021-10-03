package entity

import (
	"time"

	"github.com/n-creativesystem/rbns/domain/model"
)

type Model struct {
	ID        string `gorm:"type:varchar(256);primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) Generate() {
	m.ID = model.Generate()
}
