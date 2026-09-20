// Package config читает конфигурацию сервиса из переменных окружения.
package config

// TODO: разобрать переменные окружения из задания (шаг «Конфигурация»),
// проверить значения и вернуть ошибку, если они бессмысленны.
// Сервис с некорректной конфигурацией стартовать не должен.
type Config struct{
	HTTPAddr string
	BASE_URL string
	CODE_LENGTH int
	LOG_LEVEL string
}

const Defadr := ":8080"
const baseUrl := "http://localhost:8080"
const Deflen := 7
const baseLog := "info"


func Load() (c Config, error) {

}