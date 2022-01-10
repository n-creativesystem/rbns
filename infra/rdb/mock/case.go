package mock

import (
	"fmt"
	"testing"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/migration"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Case struct {
	db     *gorm.DB
	store  *rdb.SQLStore
	testDB TestDB
	Name   string
	Fn     func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T)
}

func (c *Case) Set(ca Case) TestCase {
	*c = ca
	return c
}

func (c *Case) Run(t *testing.T) {
	if db == nil {
		var err error
		db, err = driver.Open(c.testDB.DriverName, c.testDB.ConnStr)
		switch c.testDB.DriverName {
		case PostgreSQL:
			var count int64
			db.Table("information_schema.schemata").Where("schema_name = ?", c.testDB.Schema).Count(&count)
			if count == 0 {
				err := db.Exec(fmt.Sprintf("CREATE SCHEMA %s", c.testDB.Schema)).Error
				if err != nil {
					t.Fatal(err)
				}
			}
		case MySQL:
		}
		if err != nil {
			t.Fatal(err)
		}
		db = db.Debug()
	} else {
		c.testDB.Migration = true
	}
	c.db = db
	c.store = rdb.NewFactory(db, bus.GetBus())
	if !c.testDB.Migration {
		err := migration.MigrationTest(c.db)
		if err != nil {
			t.Fatal(err)
		}
		c.testDB.Migration = true
	}
	t.Run(c.Name, c.Fn(c.store, c.db))
	err := db.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", c.testDB.Schema)).Error
	if err != nil {
		t.Fatal(err)
	}
}

type TestCase interface {
	Run(t *testing.T)
	Set(c Case) TestCase
}

type PostgreSQLCase struct {
	schema string
	Case
}

func (pc *PostgreSQLCase) Set(c Case) TestCase {
	pc.Case = c
	return pc
}

func (pc *PostgreSQLCase) Run(t *testing.T) {
	pc.testDB = PostgreSQLTestDB(pc.schema)
	pc.Case.Run(t)
}

type MySQLCase struct {
	schema string
	Case
}

func (mc *MySQLCase) Set(c Case) TestCase {
	mc.Case = c
	return mc
}

func (mc *MySQLCase) Run(t *testing.T) {
	mc.testDB = MySQLTestDB(mc.schema)
	mc.Case.Run(t)
}

type SQLite3Case struct {
	schema string
	Case
}

func (sc *SQLite3Case) Set(c Case) TestCase {
	sc.Case = c
	return sc
}

func (sc *SQLite3Case) Run(t *testing.T) {
	sc.testDB = SQLite3TestDB()
	sc.Case.Run(t)
}

func NewCase(driverName, schema string) TestCase {
	var run TestCase
	switch driverName {
	case PostgreSQL:
		run = &PostgreSQLCase{
			schema: schema,
		}
	case MySQL:
		run = &MySQLCase{
			schema: schema,
		}
	case SQLite3:
		run = &SQLite3Case{
			schema: schema,
		}
	}
	return run
}
