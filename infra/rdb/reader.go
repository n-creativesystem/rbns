package rdb

// type reader struct {
// 	driver *gorm.DB
// }

// var _ repository.Reader = (*reader)(nil)

// func (r *reader) sess(ctx context.Context) *gorm.DB {
// 	return r.driver.Session(&gorm.Session{
// 		Context: ctx,
// 	})
// }

// func (r *reader) ApiKey(ctx context.Context) repository.ApiKey {
// 	return &apiKeyStore{
// 		db: r.sess(ctx),
// 	}
// }

// func (r *reader) Organization(ctx context.Context) repository.Organization {
// 	return &organization{
// 		db: r.sess(ctx),
// 	}
// }

// func (r *reader) User(ctx context.Context) repository.User {
// 	return &user{
// 		db: r.sess(ctx),
// 	}
// }

// func (r *reader) Resource(ctx context.Context) repository.Resource {
// 	return &resource{
// 		db: r.sess(ctx),
// 	}
// }
