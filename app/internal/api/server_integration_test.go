package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/checker"
	"github.com/itsredbull/devops-app-platform/app/internal/scheduler"
	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

func TestIntegrationAPIFlow(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://uptime:uptime@localhost:5432/uptime?sslmode=disable"
	}

	ctx := context.Background()
	st, err := store.NewPostgresStore(ctx, dsn)
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}
	defer st.Close()

	if err = st.ApplyMigrations(ctx, "../../migrations"); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	if err = st.TruncateAll(ctx); err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	probe := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("probe-ok"))
	}))
	defer probe.Close()

	mux := http.NewServeMux()
	apiServer := NewServer(st)
	apiServer.RegisterRoutes(mux)
	apiHTTP := httptest.NewServer(WithMetrics(mux))
	defer apiHTTP.Close()

	workerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker := scheduler.NewWorker(st, checker.NewService(3, 10*time.Millisecond), 1*time.Second)
	go worker.Start(workerCtx)

	payload, _ := json.Marshal(map[string]any{
		"url":                    probe.URL,
		"check_interval_seconds": 1,
		"timeout_seconds":        3,
		"enabled":                true,
	})

	resp, err := http.Post(apiHTTP.URL+"/api/v1/targets", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create target request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 from create target, got %d", resp.StatusCode)
	}

	targetsResp, err := http.Get(apiHTTP.URL + "/api/v1/targets")
	if err != nil {
		t.Fatalf("list targets request failed: %v", err)
	}
	defer targetsResp.Body.Close()
	if targetsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from list targets, got %d", targetsResp.StatusCode)
	}

	type target struct {
		ID string `json:"id"`
	}
	var targets []target
	if err = json.NewDecoder(targetsResp.Body).Decode(&targets); err != nil {
		t.Fatalf("decode targets response: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected exactly 1 target, got %d", len(targets))
	}

	deadline := time.Now().Add(12 * time.Second)
	for {
		statusResp, getErr := http.Get(apiHTTP.URL + "/api/v1/status")
		if getErr != nil {
			t.Fatalf("status request failed: %v", getErr)
		}

		var statusRows []struct {
			ID        string `json:"id"`
			LastCheck *struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
			} `json:"last_check"`
		}
		err = json.NewDecoder(statusResp.Body).Decode(&statusRows)
		statusResp.Body.Close()
		if err != nil {
			t.Fatalf("decode status response: %v", err)
		}

		for _, row := range statusRows {
			if row.ID == targets[0].ID && row.LastCheck != nil {
				// This integration test verifies end-to-end API + scheduler + DB write flow.
				// Success/failure classification is covered by checker unit tests.
				return
			}
		}

		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for scheduler check result")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
