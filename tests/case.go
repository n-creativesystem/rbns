package tests

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/n-creativesystem/rbns/tests/mock"
)

var IDs = []string{"01FA2P4T4V2FANWFBV1DFWJ8VY",
	"01FA2P4T4VEN9QKN4Y1VE3BRSB",
	"01FA2P4T4VEN9QKN4Y1Y77VQD5",
	"01FA2P4T4WHTXS335QW7JKZAEE",
	"01FA2P4T4WHTXS335QWAD8MMZ9",
	"01FA2P4T4WHTXS335QWC81148N",
	"01FA2P4T4WHTXS335QWEC6Z5XW",
	"01FA2P4T4WHTXS335QWEXWWZPS",
	"01FA2P4T4WHTXS335QWH82B8CW",
	"01FA2P4T4WHTXS335QWJS8NWKQ",
}

type Case struct {
	Name string
	Fn   func(t *testing.T)
}

type Cases []Case

func (c Cases) Run(t *testing.T) {
	for _, tt := range c {
		t.Run(tt.Name, tt.Fn)
	}
}

type MockByPostgres struct {
	Name string
	Fn   func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T)
}

type MocksByPostgres []MockByPostgres

func (cases MocksByPostgres) Run(t *testing.T) {
	for _, tt := range cases {
		db, mock := mock.NewPostgresMock()
		t.Run(tt.Name, tt.Fn(db, mock))
	}
}
