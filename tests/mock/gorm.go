package mock

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockDB struct {
	DB *gorm.DB
}

func NewPostgresMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatalln(err)
	}
	mockDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatalln(err)
	}
	return mockDB.Debug(), mock
}

func (m *MockDB) Session(ctx context.Context) *gorm.DB {
	session := &gorm.Session{
		Context:         ctx,
		CreateBatchSize: 1000,
	}
	return m.DB.Session(session)
}

func (m *MockDB) SessionSlave(ctx context.Context) *gorm.DB {
	session := &gorm.Session{
		Context:         ctx,
		CreateBatchSize: 1000,
	}
	return m.DB.Session(session)
}
