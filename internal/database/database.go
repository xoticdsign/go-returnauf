package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xoticdsign/auf-citaty/models/responses"
)

type Queuer interface {
	ListAll() []responses.Quote
	GetQuote(id string) (responses.Quote, error)
}

type Service struct {
	db *gorm.DB
}

func RunGORM() (*Service, error) {
	gormDB, err := gorm.Open(sqlite.Open(os.Getenv("DB_ADDRESS")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, gorm.ErrInvalidDB
	}
	return &Service{db: gormDB}, nil
}

func (s *Service) ListAll() []responses.Quote {
	var quotes []responses.Quote

	s.db.Table("quotes").Find(&quotes)

	return quotes
}

func (s *Service) GetQuote(id string) (responses.Quote, error) {
	var quote responses.Quote

	tx := s.db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return responses.Quote{}, gorm.ErrRecordNotFound
	}

	return quote, nil
}
