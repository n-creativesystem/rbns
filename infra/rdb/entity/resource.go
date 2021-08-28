package entity

type Resource struct {
	Method       string      `gorm:"type:varchar(256);primaryKey"`
	URI          string      `gorm:"type:varchar(256);primaryKey"`
	PermissionID string      `gorm:"type:varchar(256);primaryKey"`
	Permission   *Permission `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
