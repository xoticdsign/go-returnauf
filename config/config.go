package config

import "os"

// Структура содержащая поля с переменными окружения
type Config struct {
	ServerAddr    string
	RedisAddr     string
	RedisPassword string
	DBAddr        string
	ApiKey        string
}

// Функция подгружающая переменные окружения
func LoadConfig() Config {
	return Config{
		ServerAddr:    os.Getenv("SERVER_ADDRESS"),
		RedisAddr:     os.Getenv("REDIS_ADDRESS"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		DBAddr:        os.Getenv("DB_ADDRESS"),
		ApiKey:        os.Getenv("API_KEY"),
	}
}
