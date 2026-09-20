// Package storage хранит соответствие короткого кода и исходной ссылки.
//
// Интерфейс, ошибка и сигнатура New() — часть контракта практики.
// Менять их нельзя: против них компилируется тест, который приходит
// из пайплайна курса.
package storage

import "errors"
import "sync"

// ErrCodeTaken возвращается, если код уже занят другой ссылкой.
var ErrCodeTaken = errors.New("code already taken")

// Storage — хранилище ссылок.
type Storage interface {
	// Save сохраняет ссылку под кодом.
	// Проверить, что код свободен, и записать нужно одной неделимой операцией.
	Save(code, url string) error

	// Lookup возвращает ссылку по коду.
	Lookup(code string) (url string, ok bool)
}

	type URL_struct struct {
		urls map[string]string
		sync.RWMutex
	}
	func (u *URL_struct)  Save( code,  url string) error {
		u.Lock()
		defer u.Unlock()
		_, ok := u.urls[code]
		if ok{
			return ErrCodeTaken
		} else {
			u.urls[code] = url
			return nil
		}
	}
	func (u *URL_struct)  Lookup( code string) (string, bool) {
		u.RLock()
		defer u.RUnlock()
		value, ok := u.urls[code]
		return value, ok
	}

// New создаёт хранилище в памяти.
//
// TODO: вернуть реализацию Storage. Она должна выдерживать параллельные
// вызовы: net/http обрабатывает каждый запрос в отдельной горутине.
func New() Storage {
	m := make(map[string]string)
	u := &URL_struct{urls: m}
	return u
}
