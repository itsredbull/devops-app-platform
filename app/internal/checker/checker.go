package checker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

type Service struct {
	maxAttempts  int
	retryBackoff time.Duration
}

func NewService(maxAttempts int, retryBackoff time.Duration) *Service {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if retryBackoff < 0 {
		retryBackoff = 0
	}
	return &Service{
		maxAttempts:  maxAttempts,
		retryBackoff: retryBackoff,
	}
}

func (s *Service) RunCheck(ctx context.Context, t store.Target) store.InsertCheckInput {
	start := time.Now()
	client := &http.Client{Timeout: time.Duration(t.TimeoutSeconds) * time.Second}

	var lastErr string
	var lastCode *int

	for attempt := 1; attempt <= s.maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
		if err != nil {
			latency := int(time.Since(start).Milliseconds())
			return store.InsertCheckInput{TargetID: t.ID, LatencyMs: latency, Success: false, Error: err.Error()}
		}

		resp, err := client.Do(req)
		latency := int(time.Since(start).Milliseconds())
		if err != nil {
			lastErr = classifyErr(err)
			if attempt < s.maxAttempts {
				waitRetry(ctx, s.retryBackoff, attempt)
				continue
			}
			return store.InsertCheckInput{
				TargetID:  t.ID,
				LatencyMs: latency,
				Success:   false,
				Error:     fmt.Sprintf("request failed after %d attempt(s): %s", attempt, lastErr),
			}
		}

		code := resp.StatusCode
		_ = resp.Body.Close()
		lastCode = &code
		success := code >= 200 && code < 400
		if success {
			return store.InsertCheckInput{
				TargetID:   t.ID,
				StatusCode: lastCode,
				LatencyMs:  latency,
				Success:    true,
			}
		}

		lastErr = fmt.Sprintf("unexpected status code %d", code)
		if code >= 500 && attempt < s.maxAttempts {
			waitRetry(ctx, s.retryBackoff, attempt)
			continue
		}

		return store.InsertCheckInput{
			TargetID:   t.ID,
			StatusCode: lastCode,
			LatencyMs:  latency,
			Success:    false,
			Error:      lastErr,
		}
	}

	latency := int(time.Since(start).Milliseconds())
	return store.InsertCheckInput{
		TargetID:   t.ID,
		StatusCode: lastCode,
		LatencyMs:  latency,
		Success:    false,
		Error:      fmt.Sprintf("request failed after %d attempt(s): %s", s.maxAttempts, lastErr),
	}
}

func classifyErr(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "request timeout"
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "network timeout"
		}
		return "network error"
	}
	return err.Error()
}

func waitRetry(ctx context.Context, backoff time.Duration, attempt int) {
	if backoff <= 0 {
		return
	}
	// Linear backoff: base * attempt.
	wait := backoff * time.Duration(attempt)
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}
