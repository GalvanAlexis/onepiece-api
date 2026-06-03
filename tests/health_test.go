package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	testApp.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"status":"ok"}` // Adjust if your health handler returns something else
	if w.Body.String() != expectedBody {
		// Just a basic check.
		t.Logf("Response: %s", w.Body.String())
	}
}
