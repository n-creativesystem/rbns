package rdb

import (
	. "github.com/n-creativesystem/rbns/infra/rdb/driver/migration"
)

func (f *SQLStore) migration() error {
	return Migration(f.db)
}
