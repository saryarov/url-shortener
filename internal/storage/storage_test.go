package storage

import (
	"errors"
	"testing"
	"fmt"
	"sync"
	"sync/atomic"
)

func TestSaveAndLookup(t *testing.T) {
	s := New()

	if err := s.Save("abc", "https://example.com"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, ok := s.Lookup("abc")
	if !ok {
		t.Fatal("Lookup: код не найден")
	}
	if got != "https://example.com" {
		t.Fatalf("Lookup = %q, ожидалось %q", got, "https://example.com")
	}
}

func TestSaveRejectsTakenCode(t *testing.T) {
	s := New()

	if err := s.Save("abc", "https://first.example"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := s.Save("abc", "https://second.example"); !errors.Is(err, ErrCodeTaken) {
		t.Fatalf("Save на занятый код вернул %v, ожидалась ErrCodeTaken", err)
	}

	if got, _ := s.Lookup("abc"); got != "https://first.example" {
		t.Fatalf("занятый код перезаписан: %q", got)
	}
}

func TestLookupMissing(t *testing.T) {
	if _, ok := New().Lookup("nope"); ok {
		t.Fatal("Lookup: несуществующий код найден")
	}
}

func TestGor(t *testing.T) {
	s := New()
	var n int = 100
	code := "same-code"
	var wg sync.WaitGroup
	var okCount atomic.Int64
	var takenCount atomic.Int64
	var winnerIdx atomic.Int64
	winnerIdx.Store(-1)
	for i := 0; i < n; i++{
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := fmt.Sprintf("https://example.com/%d", i)
			err := s.Save(code, url)
			if err == nil {
				okCount.Add(1)
				winnerIdx.Store(int64(i))
			} else if errors.Is(err, ErrCodeTaken) {
				takenCount.Add(1)
			}
		}(i)
	}
	wg.Wait()
	ok := okCount.Load()
	if ok != 1{
		t.Errorf("ожидалось 1 успешное сохранение, получено %d", ok)
	}

	taken := takenCount.Load()
	if taken != int64(n-1) {
		t.Errorf("ожидалось %d ErrCodeTaken, получено %d", n-1, taken)
	}
	winner := winnerIdx.Load()
	if winner == -1 {
		t.Fatalf("ни одна горутина не сохранилась")
	}

	url, found := s.Lookup(code)
	if !found{
		t.Errorf("Lookup не нашёл код %q", code)
	}

	expected := fmt.Sprintf("https://example.com/%d", winner)
	if url != expected {
		t.Errorf("Lookup вернул %q, ожидался %q", url, expected)
	}

}


