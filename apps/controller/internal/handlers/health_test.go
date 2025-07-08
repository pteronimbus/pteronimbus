package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler_Health(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	expected := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("handler returned wrong content type: got %v want %v", ct, expected)
	}

	// Check response body
	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", response.Status)
	}

	if response.Service != "pteronimbus-controller" {
		t.Errorf("expected service 'pteronimbus-controller', got %s", response.Service)
	}

	if response.Version != "0.1.0" {
		t.Errorf("expected version '0.1.0', got %s", response.Version)
	}

	// Check timestamp is recent (within last 5 seconds)
	if time.Since(response.Timestamp) > 5*time.Second {
		t.Errorf("timestamp is too old: %v", response.Timestamp)
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "ready" {
		t.Errorf("expected status 'ready', got %s", response.Status)
	}

	if response.Service != "pteronimbus-controller" {
		t.Errorf("expected service 'pteronimbus-controller', got %s", response.Service)
	}
}

func TestHealthHandler_Live(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("GET", "/live", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "alive" {
		t.Errorf("expected status 'alive', got %s", response.Status)
	}

	if response.Service != "pteronimbus-controller" {
		t.Errorf("expected service 'pteronimbus-controller', got %s", response.Service)
	}
}

func TestHealthHandler_Healthz(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body (should be same as /health)
	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", response.Status)
	}
}

func TestHealthHandler_InvalidMethod(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Should return 405 Method Not Allowed
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHealthHandler_InvalidRoute(t *testing.T) {
	handler := NewHealthHandler()
	req, err := http.NewRequest("GET", "/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := handler.Routes()
	router.ServeHTTP(rr, req)

	// Should return 404 Not Found
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
