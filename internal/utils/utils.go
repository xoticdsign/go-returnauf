package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// Интерфейс, содержащий дополнительные методы хендлеров
type Supporter interface {
	RandInt(interval int) (int, string)
}

// Структура, реализующая Supporter
type Support struct{}

// Генерирует случайное число
func (s *Support) RandInt(count int) (int, string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := rand.Intn(count)

	id := strconv.Itoa(randInt)

	return randInt, id
}
