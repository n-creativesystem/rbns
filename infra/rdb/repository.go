package rdb

// type dbRepository struct {
// 	driver *gorm.DB
// }

// func NewRepository(driver *gorm.DB) repository.Repository {
// 	return &dbRepository{
// 		driver: driver,
// 	}
// }

// var _ repository.Repository = (*dbRepository)(nil)

// func (repo *dbRepository) NewConnection() repository.Connection {
// 	return &dbConnection{
// 		driver: repo.driver,
// 	}
// }

// func (con *dbConnection) Transaction(ctx context.Context) repository.Tx {
// 	return &transaction{
// 		db: con.driver.Session(&gorm.Session{Context: ctx}),
// 	}
// }

// type transaction struct {
// 	db *gorm.DB
// }

// var _ repository.Tx = (*transaction)(nil)

// func (t *transaction) Do(fn func(tx repository.Transaction) error) error {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			logrus.Println(err)
// 		}
// 	}()
// 	tx := t.db.Begin()
// 	defer func() {
// 		if err := recover(); err != nil {
// 			tx.Rollback()
// 			panic(err)
// 		}
// 	}()
// 	tx.SkipDefaultTransaction = true
// 	err = fn(&dbTransaction{
// 		db: tx,
// 	})
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err = tx.Commit().Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// type dbTransaction struct {
// 	db *gorm.DB
// }

// var _ repository.Transaction = (*dbTransaction)(nil)

// func (tx *dbTransaction) Permission() repository.PermissionCommand {
// 	return &permission{
// 		db: tx.db,
// 	}
// }

// func (tx *dbTransaction) Role() repository.RoleCommand {
// 	return &role{
// 		db: tx.db,
// 	}
// }

// func (tx *dbTransaction) Organization() repository.OrganizationCommand {
// 	return &organization{
// 		db: tx.db,
// 	}
// }

// func (tx *dbTransaction) User() repository.UserCommand {
// 	return &user{
// 		db: tx.db,
// 	}
// }

// func (tx *dbTransaction) Resource() repository.ResourceCommand {
// 	return &resource{
// 		db: tx.db,
// 	}
// }
