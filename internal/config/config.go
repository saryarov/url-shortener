// Package config читает конфигурацию сервиса из переменных окружения.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// TODO: разобрать переменные окружения из задания (шаг «Конфигурация»),
// проверить значения и вернуть ошибку, если они бессмысленны.
// Сервис с некорректной конфигурацией стартовать не должен.
type Config struct {
	HTTPAddr    string
	BASE_URL    string
	CODE_LENGTH int
	LOG_LEVEL   string
}

const Defadr = ":8080"
const baseUrl = "http://localhost:8080"
const Deflen = 7
const baseLog = "info"

var allowedLogLevels = map[string]bool{
	"info":  true,
	"debug": true,
	"warn":  true,
	"error": true}

func Load() (Config, error) {
	var c Config

	//====================================

	value, ok := os.LookupEnv("HTTPAddr")
	if !ok {
		c.HTTPAddr = Defadr
	} else if value == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR: пустое значение")
	} else {
		c.HTTPAddr = value
	}

	//====================================

	value, ok = os.LookupEnv("BASE_URL")
	if !ok {
		c.BASE_URL = baseUrl
	} else {
		u, err := url.Parse(value)
		if err != nil {
			return Config{}, fmt.Errorf("url.Parse: error")
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return Config{}, fmt.Errorf("BASE_URL: not https or http")
		}
		if len(u.Host) == 0 {
			return Config{}, fmt.Errorf("BASE_URL: need host")
		}
		c.BASE_URL = value
	}

	//====================================

	value, ok = os.LookupEnv("CODE_LENGTH")
	if !ok {
		c.CODE_LENGTH = Deflen
	} else {
		n, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("strconv.Atoi: error")
		}
		if n < 1 {
			return Config{}, fmt.Errorf("CODE_LENGTH: len < 1")
		} else {
			c.CODE_LENGTH = n
		}
	}

	//====================================

	value, ok = os.LookupEnv("LOG_LEVEL")
	if !ok {
		c.LOG_LEVEL = baseLog
	} else {
		if _, ok := allowedLogLevels[value]; ok {
			c.LOG_LEVEL = value
		} else {
			return Config{}, fmt.Errorf("LOG_LEVEL: only info/debug/warn/error")
		}
	}

	return c, nil
}
