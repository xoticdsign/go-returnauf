package database

import (
	"math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Quote struct {
	ID    uint   `gorm:"type:BIGINT NOT NULL PRIMARY KEY"`
	Quote string `gorm:"type:VARCHAR NOT NULL"`
}

func RunGORM() error {
	var err error

	db, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		return gorm.ErrInvalidDB
	}
	return nil
}

func ListAll() []Quote {
	var quotes []Quote

	db.Table("quotes").Find(&quotes)

	return quotes
}

func RandomQuote() Quote {
	var quote Quote

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := rand.Intn(201)

	db.Table("quotes").Where("id=?", randInt).First(&quote)

	return quote
}

func QuoteID(id string) (Quote, error) {
	var quote Quote

	tx := db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return Quote{}, gorm.ErrRecordNotFound
	}

	return quote, nil
}
