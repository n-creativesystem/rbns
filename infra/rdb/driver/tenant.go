package driver

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Tenant struct{}

func (t Tenant) Name() string {
	return "Tenant"
}

func (t Tenant) Initialize(db *gorm.DB) error {
	_ = db.Callback().Query().Before("gorm:query").Register("tenant:before_query", beforeQuery)
	return nil
}

func TenantPlugin() gorm.Plugin {
	return Tenant{}
}

type BeforeQueryInterface interface {
	BeforeQuery(tx *gorm.DB) error
}

func callMethod(db *gorm.DB, fc func(value interface{}, tx *gorm.DB) bool) {
	tx := db.Session(&gorm.Session{NewDB: true})
	if called := fc(db.Statement.ReflectValue.Interface(), tx); !called {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			db.Statement.CurDestIndex = 0
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				fc(reflect.Indirect(db.Statement.ReflectValue.Index(i)).Addr().Interface(), tx)
				db.Statement.CurDestIndex++
			}
		case reflect.Struct:
			fc(db.Statement.ReflectValue.Addr().Interface(), tx)
		}
	}
}

func beforeQuery(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		callMethod(db, func(value interface{}, tx *gorm.DB) (called bool) {
			if i, ok := value.(BeforeQueryInterface); ok {
				called = true
				db.Statement.Clauses["tenant_enabled"] = clause.Clause{}
				db.AddError(i.BeforeQuery(tx))
			}
			return called
		})
	}
}
