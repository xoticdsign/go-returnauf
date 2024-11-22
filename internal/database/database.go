package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xoticdsign/auf-citaty/internal/models/responses"
)

var db *gorm.DB

func RunGORM() error {
	dsn := os.Getenv("DB_DSN")

	var err error

	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return gorm.ErrInvalidDB
	}
	return nil
}

func ListAll() []responses.Quote {
	var quotes []responses.Quote

	db.Table("quotes").Find(&quotes)

	return quotes
}

func GetQoute(id string) (responses.Quote, error) {
	var quote responses.Quote

	tx := db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return responses.Quote{}, gorm.ErrRecordNotFound
	}

	return quote, nil
}
