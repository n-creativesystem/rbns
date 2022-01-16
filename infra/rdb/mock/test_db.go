package mock

import (
	"fmt"
	"os"

	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

const (
	PostgreSQL = driver.PostgreSQL
	MySQL      = driver.MySQL
	SQLite3    = driver.SQLite3
)

type TestDB struct {
	DriverName string
	ConnStr    string
	Schema     string
	Migration  bool
}

func SQLite3TestDB() TestDB {
	return TestDB{
		DriverName: SQLite3,
		ConnStr:    "file::memory:?cache=shared",
	}
}

func PostgreSQLTestDB(schema string) TestDB {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "api-rbac-postgres-dev"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	connStr := fmt.Sprintf("user=rbac-user password=rbac-user host=%s port=%s dbname=postgres search_path=%s sslmode=disable",
		host, port, schema)
	return TestDB{
		DriverName: PostgreSQL,
		ConnStr:    connStr,
		Schema:     schema,
	}
}

func MySQLTestDB(schema string) TestDB {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "api-rbac-mysql-dev"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	return TestDB{
		DriverName: "mysql",
		ConnStr:    fmt.Sprintf("rbac-user:rbac-user@tcp(%s:%s)/rbns_test?schema=%s", host, port, schema),
		Schema:     schema,
	}
}
