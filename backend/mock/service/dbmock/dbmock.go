package dbmock

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

func NewMockDB() *MockDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return &MockDB{
		db:   db,
		Mock: mock,
	}
}

type MockDB struct {
	postgres.Client

	db   *sql.DB
	Mock sqlmock.Sqlmock
}

func (m *MockDB) Register() {
	service.Registry[postgres.Name] = m
}

func (m *MockDB) DB() *sql.DB {
	return m.db
}

func (m *MockDB) MustMeetExpectations() {
	err := m.Mock.ExpectationsWereMet()
	if err != nil {
		panic(err)
	}
}

type AnyArg struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyArg) Match(v driver.Value) bool {
	return true
}
