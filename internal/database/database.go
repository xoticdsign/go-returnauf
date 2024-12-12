package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xoticdsign/returnauf/models/responses"
)

// Интерфейс, содержащий методы для работы с БД
type Queuer interface {
	QuotesCount() (int, error)
	ListAll() ([]responses.Quote, error)
	GetQuote(id string) (responses.Quote, error)
}

// Структура, реализующая Queuer
type DB struct {
	db *gorm.DB
}

// Запускает SQLite и возвращает структуру, реализующую Queuer
func RunGORM(dbAddr string) (*DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(dbAddr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, gorm.ErrInvalidDB
	}
	return &DB{db: gormDB}, nil
}

// Мигрирует цитаты в БД
func (d *DB) MigrateQuotes() {
	d.db.AutoMigrate(&responses.Quote{})
	d.db.Table("quotes").Create(&responses.TestQuotes)
}

// Уничтожает тестовую БД
func (d *DB) TeardownDB() {
	sqlDB, _ := d.db.DB()
	sqlDB.Close()

	os.Remove("db_test.sqlite")
}

// Возвращает количество записей в БД
func (d *DB) QuotesCount() (int, error) {
	var quotes []responses.Quote

	tx := d.db.Table("quotes").Find(&quotes)
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return int(tx.RowsAffected), nil
}

// Возвращает все записи в БД
func (d *DB) ListAll() ([]responses.Quote, error) {
	var quotes []responses.Quote

	tx := d.db.Table("quotes").Find(&quotes)
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return quotes, nil
}

// Возвращает одну запись из БД по ID
func (d *DB) GetQuote(id string) (responses.Quote, error) {
	var quote responses.Quote

	tx := d.db.Table("quotes").Where("id=?", id).First(&quote)
	if tx.RowsAffected == 0 {
		return responses.Quote{}, gorm.ErrRecordNotFound
	}
	return quote, nil
}
