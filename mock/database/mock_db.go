// mock/database/mock_db.go
package database

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockDB struct {
	DB     *sql.DB
	Mock   sqlmock.Sqlmock
	GormDB *gorm.DB
}

func NewMockDB() (*MockDB, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	return &MockDB{
		DB:     sqlDB,
		Mock:   mock,
		GormDB: gormDB,
	}, nil
}

func NewMockDBWithPing() (*MockDB, error) {
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	return &MockDB{
		DB:     sqlDB,
		Mock:   mock,
		GormDB: gormDB,
	}, nil
}

func (m *MockDB) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MockDB) ExpectPing() {
	m.Mock.ExpectPing()
}
func (m *MockDB) ExpectClose() {
	m.Mock.ExpectClose()
}
