package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xoticdsign/auf-citaty/models/responses"
)

type Database interface {
	ListAll() []responses.Quote
	GetQuote(id string) (responses.Quote, error)
}

type GormDB struct {
	db *gorm.DB
}

func RunGORM() (*GormDB, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_ADDRESS")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, gorm.ErrInvalidDB
	}
	return &GormDB{db: db}, nil
}

func (g *GormDB) ListAll() []responses.Quote {
	var quotes []responses.Quote

	g.db.Table("quotes").Find(&quotes)

	return quotes
}

func (g *GormDB) GetQuote(id string) (responses.Quote, error) {
	var quote responses.Quote

	tx := g.db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return responses.Quote{}, gorm.ErrRecordNotFound
	}

	return quote, nil
}
