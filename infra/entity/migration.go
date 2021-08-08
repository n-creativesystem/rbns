package entity

import "gorm.io/gorm"

func Migrations(db *gorm.DB) error {
	if err := db.AutoMigrate(Permission{}, Role{}, RolePermission{}, Organization{}, User{}, UserRole{}); err != nil {
		return err
	}
	return nil
}

func MigrationBack(db *gorm.DB) error {
	return db.Migrator().DropTable(Permission{}, Role{}, RolePermission{}, Organization{}, User{}, UserRole{})
}
