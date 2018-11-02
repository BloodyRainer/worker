package handlers

import (
	"fastworker/bottleneck"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRootHandler_ServeHTTP(t *testing.T) {

	// Arrange
	req, err := http.NewRequest(http.MethodGet, "local/", nil)
	if err != nil {
		t.Fatalf("could not create request %v", err)
	}

	rec := httptest.NewRecorder()
	ch := RootHandler{}

	// Act
	ch.ServeHTTP(rec, req)

	// Assert
	res := rec.Result()
	if err != nil {
		t.Fatalf("could not read response body %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %s, got %s\n", strconv.Itoa(http.StatusOK), strconv.Itoa(res.StatusCode))
	}

}

func BenchmarkBottleneckRootHandler_ServeHTTP(b *testing.B) {
	b.StopTimer()

	ch := bottleneck.Apply(RootHandler{})
	time.Sleep(10 * time.Millisecond) // wait for bottleneck to initialize

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodGet, "local/", nil)
		if err != nil {
			b.Fatalf("could not create request %v\n", err)
		}

		rec := httptest.NewRecorder()

		b.StartTimer()
		ch.ServeHTTP(rec, req)
		b.StopTimer()
	}
}
