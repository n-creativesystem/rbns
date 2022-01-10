package rdb

// type tx struct {
// 	driver *gorm.DB
// }

// var _ repository.Writer = (*tx)(nil)

// func (t *tx) Do(ctx context.Context, fn func(tx repository.Transaction) error) error {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			logrus.Println(err)
// 		}
// 	}()
// 	tx := t.driver.Session(&gorm.Session{
// 		Context: ctx,
// 	}).Begin()
// 	defer func() {
// 		if err := recover(); err != nil {
// 			tx.Rollback()
// 			panic(err)
// 		}
// 	}()
// 	tx.SkipDefaultTransaction = true
// 	err = fn(&writer{
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

// type writer struct {
// 	db *gorm.DB
// }

// var _ repository.Transaction = (*writer)(nil)

// func (tx *writer) ApiKey() repository.ApiKeyCommnad {
// 	return &apiKeyStore{
// 		db: tx.db,
// 	}
// }

// func (tx *writer) Role() repository.RoleCommand {
// 	return &role{
// 		db: tx.db,
// 	}
// }

// func (tx *writer) Organization() repository.OrganizationCommand {
// 	return &organization{
// 		db: tx.db,
// 	}
// }

// func (tx *writer) User() repository.UserCommand {
// 	return &user{
// 		db: tx.db,
// 	}
// }

// func (tx *writer) Resource() repository.ResourceCommand {
// 	return &resource{
// 		db: tx.db,
// 	}
// }
