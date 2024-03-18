package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v0/ping?name=Clark", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Endpoint /ping failed")
	}
}

func TestPintGetClientRoute(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v0/ping/getclient", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Endpoint /ping/getclient failed")
	}
}
