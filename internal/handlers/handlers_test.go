package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/logging"
	"github.com/xoticdsign/auf-citaty/internal/utils"
	"github.com/xoticdsign/auf-citaty/models/responses"
)

// Unit тесты

// Настройка Fiber для тестов
func setupTestApp(dependencies *Dependencies) *fiber.App {
	return fiber.New(fiber.Config{
		ServerHeader:  "auf-citaty",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  dependencies.Error,
		AppName:       "auf-citaty",
	})
}

// Имитация БД, реализующая методы Queuer
type MockDB struct {
	mock.Mock
}

// Имитация метода QuotesCount
func (m *MockDB) QuotesCount() (int, error) {
	args := m.Called()

	return args.Int(0), args.Error(1)
}

// Имитация метода ListAll
func (m *MockDB) ListAll() ([]responses.Quote, error) {
	args := m.Called()

	return args.Get(0).([]responses.Quote), args.Error(1)
}

// Имитация метода GetQuote
func (m *MockDB) GetQuote(id string) (responses.Quote, error) {
	args := m.Called(id)

	return args.Get(0).(responses.Quote), args.Error(1)
}

// Имитация Кэша, реализующая методы Cacher
type MockCache struct {
	mock.Mock
}

// Имитация метода Set
func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)

	return args.Error(0)
}

// Имитация метода Get
func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)

	return args.String(0), args.Error(1)
}

// Имитация Лог, реализующая методы Logger
type MockLog struct {
	mock.Mock
}

// Имитация метода Info
func (m *MockLog) Info(message string, c *fiber.Ctx) {
	m.Called(message, c)
}

// Имитация метода Warn
func (m *MockLog) Warn(message string, c *fiber.Ctx) {
	m.Called(message, c)
}

// Имитация метода Error
func (m *MockLog) Error(message string, c *fiber.Ctx) {
	m.Called(message, c)
}

// Имитация Support, содержащая дополнительные методы хендлеров
type MockSupport struct{}

// Имитация метода RandInt
func (m *MockSupport) RandInt(count int) (int, string) {
	return 1, "1"
}

// Unit тест для хендлера ListAll
func TestUnitListAll(t *testing.T) {
	cases := []struct {
		name                   string
		method                 string
		path                   string
		wantListAllToReturnErr error
		wantBodyToBe           interface{}
	}{
		{
			name:                   "general case",
			method:                 "GET",
			path:                   "/",
			wantListAllToReturnErr: nil,
			wantBodyToBe:           responses.TestQuotesForHandlers,
		},
		{
			name:                   "wrong path case",
			method:                 "GET",
			path:                   "/wrongpath",
			wantListAllToReturnErr: nil,
			wantBodyToBe:           responses.ErrDictionary[404],
		},
		{
			name:                   "wrong method case",
			method:                 "POST",
			path:                   "/",
			wantListAllToReturnErr: nil,
			wantBodyToBe:           responses.ErrDictionary[405],
		},
		{
			name:                   "empty db case",
			method:                 "GET",
			path:                   "/",
			wantListAllToReturnErr: errors.New("error"),
			wantBodyToBe:           responses.ErrDictionary[404],
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockLogger := new(MockLog)

			dependencies := &Dependencies{
				DB:     mockDB,
				Logger: mockLogger,
			}

			mockDB.On("ListAll").Return(responses.TestQuotesForHandlers, cs.wantListAllToReturnErr)

			mockLogger.On("Info", mock.Anything, mock.Anything)
			mockLogger.On("Warn", mock.Anything, mock.Anything)
			mockLogger.On("Error", mock.Anything, mock.Anything)

			mockApp := setupTestApp(dependencies)

			mockApp.Get("/", dependencies.ListAll)

			req := httptest.NewRequest(cs.method, cs.path, nil)
			res, _ := mockApp.Test(req, -1)

			gotBody, _ := io.ReadAll(res.Body)
			gotBodyStr := string(gotBody)

			wantBodyJSON, _ := json.Marshal(&cs.wantBodyToBe)
			wantBodyStr := string(wantBodyJSON)

			assert.JSONEqf(t, wantBodyStr, gotBodyStr, "got %v, while comparing output, want %v", gotBodyStr, wantBodyStr)
		})
	}
}

// Unit тест для хендлера RandomQuote
func TestUnitRandomQuote(t *testing.T) {
	cases := []struct {
		name                       string
		method                     string
		path                       string
		wantQuotesCountToReturnErr error
		wantGetQuoteToReturnErr    error
		wantCacheSetToReturnErr    error
		wantCacheGetToReturnQuote  bool
		wantCacheGetToReturnErr    error
		wantBodyToBe               interface{}
	}{
		{
			name:                       "general case",
			method:                     "GET",
			path:                       "/random",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    nil,
			wantCacheSetToReturnErr:    nil,
			wantCacheGetToReturnQuote:  false,
			wantCacheGetToReturnErr:    errors.New("error"),
			wantBodyToBe:               responses.TestQuotesForHandlers[1],
		},
		{
			name:                       "quote from cache case",
			method:                     "GET",
			path:                       "/random",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    nil,
			wantCacheSetToReturnErr:    nil,
			wantCacheGetToReturnQuote:  true,
			wantCacheGetToReturnErr:    nil,
			wantBodyToBe:               responses.TestQuotesForHandlers[1],
		},
		{
			name:                       "wrong path case",
			method:                     "GET",
			path:                       "/wrongpath",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    nil,
			wantCacheSetToReturnErr:    nil,
			wantCacheGetToReturnQuote:  false,
			wantCacheGetToReturnErr:    nil,
			wantBodyToBe:               responses.ErrDictionary[404],
		},
		{
			name:                       "wrong method case",
			method:                     "POST",
			path:                       "/random",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    nil,
			wantCacheSetToReturnErr:    nil,
			wantCacheGetToReturnQuote:  false,
			wantCacheGetToReturnErr:    nil,
			wantBodyToBe:               responses.ErrDictionary[405],
		},
		{
			name:                       "empty db case",
			method:                     "GET",
			path:                       "/random",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    errors.New("error"),
			wantCacheSetToReturnErr:    nil,
			wantCacheGetToReturnQuote:  false,
			wantCacheGetToReturnErr:    errors.New("error"),
			wantBodyToBe:               responses.ErrDictionary[404],
		},
		{
			name:                       "can't set cache case",
			method:                     "GET",
			path:                       "/random",
			wantQuotesCountToReturnErr: nil,
			wantGetQuoteToReturnErr:    nil,
			wantCacheSetToReturnErr:    errors.New("error"),
			wantCacheGetToReturnQuote:  false,
			wantCacheGetToReturnErr:    errors.New("error"),
			wantBodyToBe:               responses.TestQuotesForHandlers[1],
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockCache := new(MockCache)
			mockLogger := new(MockLog)
			mockSupport := new(MockSupport)

			dependencies := &Dependencies{
				DB:      mockDB,
				Cache:   mockCache,
				Logger:  mockLogger,
				Support: mockSupport,
			}

			mockDB.On("QuotesCount").Return(len(responses.TestQuotesForHandlers), cs.wantQuotesCountToReturnErr)
			mockDB.On("GetQuote", mock.Anything).Return(responses.TestQuotesForHandlers[1], cs.wantGetQuoteToReturnErr)

			mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(cs.wantCacheSetToReturnErr)
			mockCache.On("Get", mock.Anything).Return(responses.TestQuotesForHandlers[1].Quote, cs.wantCacheGetToReturnErr)

			mockLogger.On("Info", mock.Anything, mock.Anything)
			mockLogger.On("Warn", mock.Anything, mock.Anything)
			mockLogger.On("Error", mock.Anything, mock.Anything)

			mockApp := setupTestApp(dependencies)

			mockApp.Get("/random", dependencies.RandomQuote)

			req := httptest.NewRequest(cs.method, cs.path, nil)
			resp, _ := mockApp.Test(req, -1)

			gotBody, _ := io.ReadAll(resp.Body)
			gotBodyStr := string(gotBody)

			wantBodyJSON, _ := json.Marshal(&cs.wantBodyToBe)
			wantBodyStr := string(wantBodyJSON)

			assert.JSONEqf(t, wantBodyStr, gotBodyStr, "got %v, while comparing output, want %v", gotBodyStr, wantBodyStr)
		})
	}
}

// Unit тест для хендлера QuoteID
func TestUnitQuoteID(t *testing.T) {
	cases := []struct {
		name                      string
		method                    string
		path                      string
		wantGetQuoteToReturnErr   error
		wantCacheSetToReturnErr   error
		wantCacheGetToReturnQuote bool
		wantCacheGetToReturnErr   error
		wantBodyToBe              interface{}
	}{
		{
			name:                      "general case",
			method:                    "GET",
			path:                      "/1",
			wantGetQuoteToReturnErr:   nil,
			wantCacheSetToReturnErr:   nil,
			wantCacheGetToReturnQuote: false,
			wantCacheGetToReturnErr:   errors.New("error"),
			wantBodyToBe:              responses.TestQuotesForHandlers[1],
		},
		{
			name:                      "quote from cache case",
			method:                    "GET",
			path:                      "/1",
			wantGetQuoteToReturnErr:   nil,
			wantCacheSetToReturnErr:   nil,
			wantCacheGetToReturnQuote: true,
			wantCacheGetToReturnErr:   nil,
			wantBodyToBe:              responses.TestQuotesForHandlers[1],
		},
		{
			name:                      "wrong path case",
			method:                    "GET",
			path:                      "/wrongpath",
			wantGetQuoteToReturnErr:   nil,
			wantCacheSetToReturnErr:   nil,
			wantCacheGetToReturnQuote: false,
			wantCacheGetToReturnErr:   nil,
			wantBodyToBe:              responses.ErrDictionary[404],
		},
		{
			name:                      "wrong method case",
			method:                    "POST",
			path:                      "/1",
			wantGetQuoteToReturnErr:   nil,
			wantCacheSetToReturnErr:   nil,
			wantCacheGetToReturnQuote: false,
			wantCacheGetToReturnErr:   nil,
			wantBodyToBe:              responses.ErrDictionary[405],
		},
		{
			name:                      "empty db case",
			method:                    "GET",
			path:                      "/1",
			wantGetQuoteToReturnErr:   errors.New("error"),
			wantCacheSetToReturnErr:   nil,
			wantCacheGetToReturnQuote: false,
			wantCacheGetToReturnErr:   errors.New("error"),
			wantBodyToBe:              responses.ErrDictionary[404],
		},
		{
			name:                      "can't set cache case",
			method:                    "GET",
			path:                      "/1",
			wantGetQuoteToReturnErr:   nil,
			wantCacheSetToReturnErr:   errors.New("error"),
			wantCacheGetToReturnQuote: false,
			wantCacheGetToReturnErr:   errors.New("error"),
			wantBodyToBe:              responses.TestQuotesForHandlers[1],
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockCache := new(MockCache)
			mockLogger := new(MockLog)

			dependencies := &Dependencies{
				DB:     mockDB,
				Cache:  mockCache,
				Logger: mockLogger,
			}

			mockDB.On("GetQuote", mock.Anything).Return(responses.TestQuotesForHandlers[1], cs.wantGetQuoteToReturnErr)

			mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(cs.wantCacheSetToReturnErr)
			mockCache.On("Get", mock.Anything).Return(responses.TestQuotesForHandlers[1].Quote, cs.wantCacheGetToReturnErr)

			mockLogger.On("Info", mock.Anything, mock.Anything)
			mockLogger.On("Warn", mock.Anything, mock.Anything)
			mockLogger.On("Error", mock.Anything, mock.Anything)

			mockApp := setupTestApp(dependencies)

			mockApp.Get("/:id", dependencies.QuoteID)

			req := httptest.NewRequest(cs.method, cs.path, nil)
			resp, _ := mockApp.Test(req, -1)

			gotBody, _ := io.ReadAll(resp.Body)
			gotBodyStr := string(gotBody)

			wantBodyJSON, _ := json.Marshal(&cs.wantBodyToBe)
			wantBodyStr := string(wantBodyJSON)

			assert.JSONEqf(t, wantBodyStr, gotBodyStr, "got %v, while comparing output, want %v", gotBodyStr, wantBodyStr)
		})
	}
}

// Integration тесты

// Настройка БД для интеграционных тестов
func setupTestDB(emptyDB bool) *database.DB {
	DB, _ := database.RunGORM("db_test.sqlite")

	if !emptyDB {
		DB.MigrateQuotes()
	}

	return DB
}

// Настройка Кэша для интеграционных тестов
func setupTestCache(emptyCache bool) *cache.Cache {
	Cache, _ := cache.RunRedis("127.0.0.1:6379", "")

	if !emptyCache {
		Cache.PopulateCache()
	}

	return Cache
}

// Integration тест для хендлера ListAll
func TestIntegrationListAll(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		path         string
		emptyDB      bool
		emptyCache   bool
		wantStatus   int
		wantBodyToBe interface{}
	}{
		{
			name:         "regular case",
			method:       "GET",
			path:         "/",
			emptyDB:      false,
			emptyCache:   true,
			wantStatus:   200,
			wantBodyToBe: responses.TestQuotes,
		},
		{
			name:         "wrong method case",
			method:       "POST",
			path:         "/",
			emptyDB:      false,
			emptyCache:   true,
			wantStatus:   405,
			wantBodyToBe: responses.ErrDictionary[405],
		},
		{
			name:         "wrong path case",
			method:       "GET",
			path:         "/wrongpath",
			emptyDB:      false,
			emptyCache:   true,
			wantStatus:   404,
			wantBodyToBe: responses.ErrDictionary[404],
		},
		{
			name:         "empty db case",
			method:       "GET",
			path:         "/",
			emptyDB:      true,
			emptyCache:   true,
			wantStatus:   404,
			wantBodyToBe: responses.ErrDictionary[404],
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			DB := setupTestDB(cs.emptyDB)
			defer DB.TeardownDB()

			Cache := setupTestCache(cs.emptyCache)
			defer Cache.TeardownCache()

			Log, _ := logging.RunZap()

			dependencies := &Dependencies{
				DB:      DB,
				Cache:   Cache,
				Logger:  Log,
				Support: &utils.Support{},
			}

			testApp := setupTestApp(dependencies)

			testApp.Use(requestid.New(requestid.Config{
				Generator:  uuid.NewString,
				ContextKey: "uuid",
			}))

			testApp.Get("/", dependencies.ListAll)

			req := httptest.NewRequest(cs.method, cs.path, nil)
			resp, _ := testApp.Test(req, -1)

			assert.Equalf(t, cs.wantStatus, resp.StatusCode, "got %v, while comparing returned status, want %v", resp.StatusCode, cs.wantStatus)

			gotBody, _ := io.ReadAll(resp.Body)
			gotBodyStr := string(gotBody)

			wantBodyJSON, _ := json.Marshal(&cs.wantBodyToBe)
			wantBodyStr := string(wantBodyJSON)

			assert.JSONEqf(t, wantBodyStr, gotBodyStr, "got %v, while comparing output, want %v", gotBodyStr, wantBodyStr)
		})
	}
}
