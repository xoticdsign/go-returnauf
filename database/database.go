package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Quote struct {
	ID    int    `gorm:"type:BIGINT NOT NULL PRIMARY KEY"`
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

func GetQoute(id string) (Quote, error) {
	var quote Quote

	tx := db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return Quote{}, gorm.ErrRecordNotFound
	}

	return quote, nil
}
