package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitHealthySucceedsImmediately(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := WaitHealthy(srv.URL, time.Second, 10*time.Millisecond); err != nil {
		t.Fatalf("WaitHealthy: %v", err)
	}
}

func TestWaitHealthyRetriesThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := WaitHealthy(srv.URL, time.Second, 5*time.Millisecond); err != nil {
		t.Fatalf("WaitHealthy: %v", err)
	}
	if calls < 3 {
		t.Errorf("got %d calls, want at least 3", calls)
	}
}

func TestWaitHealthyTimesOut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	err := WaitHealthy(srv.URL, 30*time.Millisecond, 10*time.Millisecond)
	if err == nil {
		t.Fatal("WaitHealthy: got nil error, want a timeout error")
	}
}
