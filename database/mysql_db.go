package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL(user, password, host, port, dbName string) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Please check the environment variables")
		log.Fatal(err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &MySQL{db: db}, nil
}

func (m *MySQL) GetDB() *gorm.DB {
	return m.db
}

func (m *MySQL) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("database connection closed")
	return nil
}
