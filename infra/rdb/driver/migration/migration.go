package migration

import (
	"github.com/n-creativesystem/rbns/infra/entity"
	"gorm.io/gorm"
)

var (
	migrationTable = []interface{}{
		entity.ApiKey{}, entity.Permission{}, entity.Role{},
		entity.Organization{}, entity.User{}, entity.Tenant{}, entity.LoginUser{},
	}
)

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(migrationTable...)
}

func MigrationTest(db *gorm.DB) error {
	for _, table := range migrationTable {
		if db.Migrator().HasTable(table) {
			db.Delete(&table, "1 = 1")
		}
	}
	return Migration(db)
}
