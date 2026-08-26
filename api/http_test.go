package api

import (
	"drupal-scheduler/service"
	"drupal-scheduler/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealth(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer db.Close()
	s := service.New(db)
	defer s.Close()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	New(s).Handler().ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
