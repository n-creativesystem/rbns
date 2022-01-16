package plugins

import (
	"context"

	"github.com/n-creativesystem/rbns/ncsfw/tenants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Tenant string

func (t Tenant) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	value := string(t)
	if value == "" {
		value = tenants.FromTenantContext(ctx)
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{value},
	}
}

// func (t Tenant) QueryClauses(f *schema.Field) []clause.Interface {
// 	return []clause.Interface{TenantQueryClauses{Field: f, Value: string(t)}}
// }

// func (t Tenant) UpdateClauses(f *schema.Field) []clause.Interface {
// 	return []clause.Interface{TenantQueryClauses{Field: f, Value: string(t)}}
// }

// func (t Tenant) CreateClauses(f *schema.Field) []clause.Interface {
// 	return []clause.Interface{TenantValuesCaluses{Field: f, Value: string(t)}}
// }

// func (t Tenant) DeleteClauses(f *schema.Field) []clause.Interface {
// 	return []clause.Interface{TenantQueryClauses{Field: f, Value: string(t)}}
// }

// func (t Tenant) String() string {
// 	return string(t)
// }

// type TenantQueryClauses struct {
// 	Field *schema.Field
// 	Value string
// }

// func (TenantQueryClauses) Name() string { return "" }

// func (TenantQueryClauses) Build(clause.Builder) {}

// func (TenantQueryClauses) MergeClause(*clause.Clause) {}

// func (t TenantQueryClauses) ModifyStatement(stmt *gorm.Statement) {
// 	if _, ok := stmt.Clauses["tenant_enabled"]; !ok {
// 		value := t.Value
// 		if value == "" {
// 			value = contexts.FromTenantContext(stmt.Context)
// 		}
// 		if c, ok := stmt.Clauses["WHERE"]; ok {
// 			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
// 				for _, expr := range where.Exprs {
// 					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
// 						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
// 						c.Expression = where
// 						stmt.Clauses["WHERE"] = c
// 						break
// 					}
// 				}
// 			}
// 		}
// 		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
// 			clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: t.Field.DBName}, Value: value},
// 		}})
// 		stmt.Clauses["tenant_enabled"] = clause.Clause{}
// 	}
// }

// type TenantValuesCaluses struct {
// 	Field *schema.Field
// 	Value string
// }

// func (TenantValuesCaluses) Name() string { return "" }

// func (TenantValuesCaluses) Build(clause.Builder) {}

// func (TenantValuesCaluses) MergeClause(*clause.Clause) {}

// func (t TenantValuesCaluses) ModifyStatement(stmt *gorm.Statement) {
// 	tenant := contexts.FromTenantContext(stmt.Context)
// 	values := clause.Values{Columns: []clause.Column{
// 		{
// 			Name: t.Field.DBName,
// 		},
// 	}}
// 	switch stmt.ReflectValue.Kind() {
// 	case reflect.Slice, reflect.Array:
// 		l := stmt.ReflectValue.Len()
// 		values.Values = make([][]interface{}, stmt.ReflectValue.Len())
// 		for i := 0; i < l; i++ {
// 			rv := reflect.Indirect(stmt.ReflectValue.Index(i))
// 			values.Values[i] = make([]interface{}, 1)
// 			field := stmt.Schema.FieldsByDBName[t.Field.DBName]
// 			if t.Value == "" {
// 				_ = field.Set(rv, tenant)
// 				values.Values[i][0], _ = field.ValueOf(rv)
// 			} else {
// 				values.Values[i][0] = t.Value
// 			}
// 		}
// 	case reflect.Struct:
// 		values.Values = [][]interface{}{make([]interface{}, 1)}
// 		field := stmt.Schema.FieldsByDBName[t.Field.DBName]
// 		if t.Value == "" {
// 			_ = field.Set(stmt.ReflectValue, tenant)
// 			values.Values[0][0], _ = field.ValueOf(stmt.ReflectValue)
// 		} else {
// 			values.Values[0][0] = t.Value
// 		}
// 	}
// 	stmt.AddClause(values)
// }
