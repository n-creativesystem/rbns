package infra

import (
	"os"
	"path/filepath"

	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	"github.com/n-creativesystem/rbns/storage"
	"github.com/spf13/viper"
)

const (
	postgresStorageType = "postgres"
	mysqlStorageType    = "mysql"
	inMemoryStorageType = "inmemory"
)

func NewFactory(typ string) (storage.Factory, map[string]interface{}, error) {
	v := viper.New()
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
	v.AddConfigPath(filepath.Join(homedir, ".rbns"))
	v.AddConfigPath("/etc/rbns")
	v.AddConfigPath("./")
	v.SetConfigName("storage")
	if err := v.ReadInConfig(); err != nil {
		return nil, nil, err
	}

	var factory storage.Factory
	switch typ {
	case postgresStorageType:
		dsn := v.GetString("DSN")
		dsn = os.ExpandEnv(dsn)
		db, _ := postgres.New(dsn)
		factory = rdb.NewFactory(db)
	case mysqlStorageType:
		dsn := v.GetString("DSN")
		dsn = os.ExpandEnv(dsn)
		db, _ := mysql.New(dsn)
		factory = rdb.NewFactory(db)
	case inMemoryStorageType:
	}
	mp := v.AllSettings()
	return factory, mp, nil
}
