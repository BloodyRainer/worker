package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRootHandler_ServeHTTP(t *testing.T) {

	// Arrange
	req, err := http.NewRequest("GET", "local/", nil)
	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	rec := httptest.NewRecorder()
	rh := RootHandler{}

	// Act
	rh.ServeHTTP(rec, req)

	// Assert
	res := rec.Result()
	if err != nil {
		t.Fatalf("could not read response body %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %s, got %s\n", strconv.Itoa(http.StatusOK), strconv.Itoa(res.StatusCode))
	}

}
