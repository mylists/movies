package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mylists/movies"
	"github.com/mylists/mtvl/pkg/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStandaloneMoviesServerEndpoints(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	mod := movies.NewModule(db)
	authProvider := auth.NewJWTAuthProvider(db, "test-secret")
	r := buildRouter(mod, authProvider)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	req = httptest.NewRequest("GET", "/api/v1/movies", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestClientIPMiddleware(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	mod := movies.NewModule(db)
	authProvider := auth.NewJWTAuthProvider(db, "test-secret")
	r := buildRouter(mod, authProvider)

	var recordedIP string
	r.Get("/test-client-ip", func(w http.ResponseWriter, req *http.Request) {
		recordedIP = middleware.GetClientIP(req.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test-client-ip", nil)
	req.RemoteAddr = "198.51.100.42:1234"
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if recordedIP != "198.51.100.42" {
		t.Errorf("expected client IP 198.51.100.42, got %q", recordedIP)
	}
	if req.RemoteAddr != "198.51.100.42:1234" {
		t.Errorf("expected RemoteAddr to remain unmutated, got %q", req.RemoteAddr)
	}
}

