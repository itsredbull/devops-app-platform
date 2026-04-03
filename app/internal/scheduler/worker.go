package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/checker"
	"github.com/itsredbull/devops-app-platform/app/internal/metrics"
	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

type Worker struct {
	store    *store.PostgresStore
	checker  *checker.Service
	interval time.Duration

	mu       sync.Mutex
	lastRun  map[string]time.Time
	inFlight map[string]bool
}

func NewWorker(st *store.PostgresStore, ch *checker.Service, interval time.Duration) *Worker {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	return &Worker{
		store:    st,
		checker:  ch,
		interval: interval,
		lastRun:  make(map[string]time.Time),
		inFlight: make(map[string]bool),
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("scheduler started with interval=%s", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	targets, err := w.store.ListEnabledTargets(ctx)
	if err != nil {
		log.Printf("scheduler: list targets failed: %v", err)
		return
	}

	metrics.SetEnabledTargets(len(targets))
	now := time.Now()

	for _, target := range targets {
		t := target
		if t.CheckIntervalSeconds <= 0 {
			t.CheckIntervalSeconds = 30
		}

		w.mu.Lock()
		last, hasLast := w.lastRun[t.ID]
		if w.inFlight[t.ID] {
			w.mu.Unlock()
			continue
		}
		if hasLast && now.Sub(last) < time.Duration(t.CheckIntervalSeconds)*time.Second {
			w.mu.Unlock()
			continue
		}
		w.inFlight[t.ID] = true
		w.lastRun[t.ID] = now
		w.mu.Unlock()

		go w.runSingle(ctx, t)
	}
}

func (w *Worker) runSingle(ctx context.Context, t store.Target) {
	result := w.checker.RunCheck(ctx, t)
	if err := w.store.InsertCheck(ctx, result); err != nil {
		log.Printf("scheduler: insert check failed target=%s err=%v", t.ID, err)
	}

	statusCode := 0
	if result.StatusCode != nil {
		statusCode = *result.StatusCode
	}
	metrics.RecordCheck(result.Success, statusCode, float64(result.LatencyMs))

	w.mu.Lock()
	delete(w.inFlight, t.ID)
	w.mu.Unlock()
}
