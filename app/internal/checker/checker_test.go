package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

func TestRunCheckSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	svc := NewService(3, 10*time.Millisecond)
	out := svc.RunCheck(context.Background(), store.Target{
		ID:             "t1",
		URL:            srv.URL,
		TimeoutSeconds: 2,
	})

	if !out.Success {
		t.Fatalf("expected success=true, got false with error=%q", out.Error)
	}
	if out.StatusCode == nil || *out.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %#v", out.StatusCode)
	}
}

func TestRunCheckRetriesThenSuccess(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	svc := NewService(3, 10*time.Millisecond)
	out := svc.RunCheck(context.Background(), store.Target{
		ID:             "t2",
		URL:            srv.URL,
		TimeoutSeconds: 2,
	})

	if !out.Success {
		t.Fatalf("expected success after retries, got error=%q", out.Error)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
}

func TestRunCheckTimeoutFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	svc := NewService(2, 5*time.Millisecond)
	out := svc.RunCheck(context.Background(), store.Target{
		ID:             "t3",
		URL:            srv.URL,
		TimeoutSeconds: 1,
	})

	if out.Success {
		t.Fatal("expected timeout failure")
	}
	if out.Error == "" {
		t.Fatal("expected non-empty error on timeout")
	}
}
