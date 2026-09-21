// Package httpapi содержит HTTP-слой сервиса: маршруты, обработчики и middleware.
package httpapi

// TODO: собрать два маршрута из задания (шаг «HTTP»), реализовать обработчики
// и запись в логи.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"urlshortener/internal/config"
	"urlshortener/internal/storage"
)

type Handler struct {
	storage storage.Storage
	cfg     config.Config
	log     *slog.Logger
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	Error string `json:"error"`
}

type createLinkRequest struct {
	URL string `json:"url"`
}

type createLinkResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func generateCode(n int) string {
	alph := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	b := make([]byte, n)
	for index := range b {
		b[index] = alph[rand.IntN(len(alph))]
	}
	return string(b)
}

func validUrl(s string) error {
	if len(s) < 1 {
		return fmt.Errorf("length < 1")
	}
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("error parse url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("need Scheme https or http")
	}
	if u.Host == "" {
		return fmt.Errorf("отсутствует хост")
	}
	return nil
}

func buildShortUrl(base, code string) string {
	base = strings.TrimSuffix(base, "/")
	return base + "/r/" + code
}

func (h *Handler) createLink(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		writeError(w, 400, "Content type must be application/json")
		return
	}
	var req createLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	err = validUrl(req.URL)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var savedCode string
	for i := 0; i < 5; i++ {
		code := generateCode(h.cfg.CODE_LENGTH)
		err = h.storage.Save(code, req.URL)
		if err == nil {
			savedCode = code
			break
		}
		if !errors.Is(err, storage.ErrCodeTaken) {
			writeError(w, 500, err.Error())
			return
		}
	}
	if savedCode == "" {
		writeError(w, 500, "не удалось подобрать свободный код")
		return
	}
	short := buildShortUrl(h.cfg.BASE_URL, savedCode)
	writeJSON(w, 201, createLinkResponse{Code: savedCode, ShortURL: short})
}
