package service

import (
	"errors"
	"testing"
)

type testStore struct {
	data               map[string]string
	visits             map[string]int64
	incrementCalls     map[string]int
	saveErr            error
	loadErr            error
	getVisitsErr       error
	incrementVisitsErr error
}

func newTestStore() *testStore {
	return &testStore{
		data:           make(map[string]string),
		visits:         make(map[string]int64),
		incrementCalls: make(map[string]int),
	}
}

func (s *testStore) Save(shortCode string, originalURL string) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.data[shortCode] = originalURL
	return nil
}

func (s *testStore) Load(shortCode string) (string, error) {
	if s.loadErr != nil {
		return "", s.loadErr
	}
	v, ok := s.data[shortCode]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (s *testStore) IncrementVisits(shortCode string) error {
	if s.incrementVisitsErr != nil {
		return s.incrementVisitsErr
	}
	s.incrementCalls[shortCode]++
	s.visits[shortCode]++
	return nil
}

func (s *testStore) GetVisits(shortCode string) (int64, error) {
	if s.getVisitsErr != nil {
		return 0, s.getVisitsErr
	}
	return s.visits[shortCode], nil
}

func (s *testStore) Delete(shortCode string) error {
	delete(s.data, shortCode)
	delete(s.visits, shortCode)
	delete(s.incrementCalls, shortCode)
	return nil
}

func TestShorten(t *testing.T) {
	st := newTestStore()
	svc := NewShortenerService(st)

	shortCode, err := svc.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("Shorten returned error: %v", err)
	}
	if shortCode == "" {
		t.Fatal("shortCode should not be empty")
	}
	if got := st.data[shortCode]; got != "https://example.com" {
		t.Fatalf("stored url mismatch, got %q", got)
	}
}

func TestGetOriginalURL(t *testing.T) {
	st := newTestStore()
	st.data["abc123"] = "https://example.com/path"
	svc := NewShortenerService(st)

	got, err := svc.GetOriginalURL("abc123")
	if err != nil {
		t.Fatalf("GetOriginalURL returned error: %v", err)
	}
	if got != "https://example.com/path" {
		t.Fatalf("url mismatch, got %q", got)
	}
	if st.incrementCalls["abc123"] != 1 {
		t.Fatalf("IncrementVisits should be called once, got %d", st.incrementCalls["abc123"])
	}
}

func TestGetStats(t *testing.T) {
	st := newTestStore()
	st.visits["xyz"] = 7
	svc := NewShortenerService(st)

	got, err := svc.GetStats("xyz")
	if err != nil {
		t.Fatalf("GetStats returned error: %v", err)
	}
	if got != 7 {
		t.Fatalf("visits mismatch, got %d", got)
	}
}
