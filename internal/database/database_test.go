package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xoticdsign/returnauf/models/responses"
)

// Настройка GORM для тестов
func setupTestDB(emptyDB bool) *DB {
	DB, _ := RunGORM("db_test.sqlite")

	if !emptyDB {
		DB.MigrateQuotes()
	}

	return DB
}

// Unit тест для функции RunGORM
func TestUnitRunGORM(t *testing.T) {
	cases := []struct {
		name                   string
		wantRunGORMToReturnErr error
		wantLoggerToBe         logger.Interface
	}{
		{
			name:                   "general case",
			wantRunGORMToReturnErr: nil,
			wantLoggerToBe:         logger.Default.LogMode(logger.Silent),
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			gotDB, gotErr := RunGORM("db_test.sqlite")
			if gotErr != nil {
				assert.Equal(t, cs.wantRunGORMToReturnErr, gotErr)
			} else {
				gotLogger := gotDB.db.Config.Logger
				sqlDB, _ := gotDB.db.DB()
				defer sqlDB.Close()

				assert.Equal(t, cs.wantLoggerToBe, gotLogger)
			}

			defer os.Remove("db_test.sqlite")
		})
	}
}

// Unit тест для функции QuotesCount
func TestUnitQuotesCount(t *testing.T) {
	cases := []struct {
		name                         string
		emptyDB                      bool
		wantQuotesCountToReturnCount int
		wantQuotesCountToReturnErr   error
	}{
		{
			name:                         "general case",
			emptyDB:                      false,
			wantQuotesCountToReturnCount: len(responses.TestQuotes),
			wantQuotesCountToReturnErr:   nil,
		},
		{
			name:                         "empty db case",
			emptyDB:                      true,
			wantQuotesCountToReturnCount: 0,
			wantQuotesCountToReturnErr:   gorm.ErrRecordNotFound,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			DB := setupTestDB(cs.emptyDB)
			defer DB.TeardownDB()

			gotCount, gotErr := DB.QuotesCount()
			if gotErr != nil {
				assert.Equal(t, cs.wantQuotesCountToReturnErr, gotErr)
			} else {
				assert.Equal(t, cs.wantQuotesCountToReturnCount, gotCount)
			}
		})
	}
}

// Unit тест для функции ListAll
func TestUnitListAll(t *testing.T) {
	cases := []struct {
		name                      string
		emptyDB                   bool
		wantListAllToReturnQuotes []responses.Quote
		wantListAllToReturnErr    error
	}{
		{
			name:                      "general case",
			emptyDB:                   false,
			wantListAllToReturnQuotes: responses.TestQuotes,
			wantListAllToReturnErr:    nil,
		},
		{
			name:                      "empty db case",
			emptyDB:                   true,
			wantListAllToReturnQuotes: nil,
			wantListAllToReturnErr:    gorm.ErrRecordNotFound,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			DB := setupTestDB(cs.emptyDB)
			defer DB.TeardownDB()

			gotQuotes, gotErr := DB.ListAll()
			if gotErr != nil {
				assert.Equal(t, cs.wantListAllToReturnErr, gotErr)
			} else {
				assert.Equal(t, cs.wantListAllToReturnQuotes, gotQuotes)
			}
		})
	}
}

// Unit тест для функции GetQuote
func TestUnitGetQuote(t *testing.T) {
	cases := []struct {
		name                      string
		input                     string
		emptyDB                   bool
		wantGetQuoteToReturnQuote responses.Quote
		wantGetQuoteToReturnErr   error
	}{
		{
			name:                      "general case",
			input:                     "1",
			emptyDB:                   false,
			wantGetQuoteToReturnQuote: responses.TestQuotes[0],
			wantGetQuoteToReturnErr:   nil,
		},
		{
			name:                      "empty db case",
			emptyDB:                   true,
			wantGetQuoteToReturnQuote: responses.Quote{},
			wantGetQuoteToReturnErr:   gorm.ErrRecordNotFound,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			DB := setupTestDB(cs.emptyDB)
			defer DB.TeardownDB()

			gotQuote, gotErr := DB.GetQuote(cs.input)
			if gotErr != nil {
				assert.Equal(t, cs.wantGetQuoteToReturnErr, gotErr)
			} else {
				assert.Equal(t, cs.wantGetQuoteToReturnQuote, gotQuote)
			}
		})
	}
}
