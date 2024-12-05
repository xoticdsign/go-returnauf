package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoticdsign/auf-citaty/models/responses"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Массив с цитатами для тестов
var quotes = []responses.Quote{
	{ID: 1, Quote: "Mock quote 1"},
	{ID: 2, Quote: "Mock quote 2"},
	{ID: 3, Quote: "Mock quote 3"},
}

// Настройка GORM для тестов
func setup(emptyDB bool) *DB {
	DB, _ := RunGORM("db_test.sqlite")

	if emptyDB == false {
		DB.db.AutoMigrate(&responses.Quote{})
		DB.db.Table("quotes").Create(&quotes)
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
			sqlDB, _ := gotDB.db.DB()
			defer sqlDB.Close()

			gotLogger := gotDB.db.Config.Logger

			assert.Equalf(t, cs.wantLoggerToBe, gotLogger, "got %v, while comparing logger from the config, want %v", gotLogger, cs.wantLoggerToBe)
			assert.Equalf(t, cs.wantRunGORMToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantRunGORMToReturnErr)

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
			wantQuotesCountToReturnCount: len(quotes),
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
			DB := setup(cs.emptyDB)

			sqlDB, _ := DB.db.DB()
			defer sqlDB.Close()

			gotCount, gotErr := DB.QuotesCount()

			assert.Equalf(t, cs.wantQuotesCountToReturnCount, gotCount, "got %v, while comparing returned count, want %v", gotCount, cs.wantQuotesCountToReturnCount)
			assert.Equalf(t, cs.wantQuotesCountToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantQuotesCountToReturnErr)

			defer os.Remove("db_test.sqlite")
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
			wantListAllToReturnQuotes: quotes,
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
			DB := setup(cs.emptyDB)

			sqlDB, _ := DB.db.DB()
			defer sqlDB.Close()

			gotQuotes, gotErr := DB.ListAll()

			assert.Equalf(t, cs.wantListAllToReturnQuotes, gotQuotes, "got %v, while comparing returned quotes, want %v", gotQuotes, cs.wantListAllToReturnQuotes)
			assert.Equalf(t, cs.wantListAllToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantListAllToReturnErr)

			defer os.Remove("db_test.sqlite")
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
			wantGetQuoteToReturnQuote: quotes[0],
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
			DB := setup(cs.emptyDB)

			sqlDB, _ := DB.db.DB()
			defer sqlDB.Close()

			gotQuote, gotErr := DB.GetQuote(cs.input)

			assert.Equalf(t, cs.wantGetQuoteToReturnQuote, gotQuote, "got %v, while comparing returned quote, want %v", gotQuote, cs.wantGetQuoteToReturnQuote)
			assert.Equalf(t, cs.wantGetQuoteToReturnErr, gotErr, "got %v, while comparing returned error, want %v", gotErr, cs.wantGetQuoteToReturnErr)

			defer os.Remove("db_test.sqlite")
		})
	}
}
