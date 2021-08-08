package infra

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func init() {
	di.MustRegister(NewRepository)
}

type dbRepository struct {
	driver dao.DataBase
}

func NewRepository(driver dao.DataBase) repository.Repository {
	return &dbRepository{
		driver: driver,
	}
}

var _ repository.Repository = (*dbRepository)(nil)

func (repo *dbRepository) NewConnection() repository.Connection {
	return &dbConnection{
		driver: repo.driver,
	}
}

type dbConnection struct {
	driver dao.DataBase
}

var _ repository.Connection = (*dbConnection)(nil)

func (con *dbConnection) Permission(ctx context.Context) repository.Permission {
	return &permission{
		db: con.driver.Session(ctx),
	}
}

func (con *dbConnection) Role(ctx context.Context) repository.Role {
	return &role{
		db: con.driver.Session(ctx),
	}
}

func (con *dbConnection) Organization(ctx context.Context) repository.Organization {
	return &organization{
		db: con.driver.Session(ctx),
	}
}

func (con *dbConnection) User(ctx context.Context) repository.User {
	return &user{
		db: con.driver.Session(ctx),
	}
}

func (con *dbConnection) Resource(ctx context.Context) repository.Resource {
	return &resource{
		db: con.driver.Session(ctx),
	}
}

func (con *dbConnection) Transaction(ctx context.Context) repository.Tx {
	return &transaction{
		db: con.driver.Session(ctx),
	}
}

type transaction struct {
	db *gorm.DB
}

var _ repository.Tx = (*transaction)(nil)

func (t *transaction) Do(fn func(tx repository.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Println(err)
		}
	}()
	tx := t.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	tx.SkipDefaultTransaction = true
	err = fn(&dbTransaction{
		db: tx,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

type dbTransaction struct {
	db *gorm.DB
}

var _ repository.Transaction = (*dbTransaction)(nil)

func (tx *dbTransaction) Permission() repository.PermissionCommand {
	return &permission{
		db: tx.db,
	}
}

func (tx *dbTransaction) Role() repository.RoleCommand {
	return &role{
		db: tx.db,
	}
}

func (tx *dbTransaction) Organization() repository.OrganizationCommand {
	return &organization{
		db: tx.db,
	}
}

func (tx *dbTransaction) User() repository.UserCommand {
	return &user{
		db: tx.db,
	}
}

func (tx *dbTransaction) Resource() repository.ResourceCommand {
	return &resource{
		db: tx.db,
	}
}
