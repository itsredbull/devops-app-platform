package store

import "time"

type Target struct {
	ID                   string    `json:"id"`
	URL                  string    `json:"url"`
	CheckIntervalSeconds int       `json:"check_interval_seconds"`
	TimeoutSeconds       int       `json:"timeout_seconds"`
	Enabled              bool      `json:"enabled"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateTargetInput struct {
	URL                  string `json:"url"`
	CheckIntervalSeconds int    `json:"check_interval_seconds"`
	TimeoutSeconds       int    `json:"timeout_seconds"`
	Enabled              bool   `json:"enabled"`
}

type Check struct {
	ID         string    `json:"id"`
	TargetID   string    `json:"target_id"`
	CheckedAt  time.Time `json:"checked_at"`
	StatusCode *int      `json:"status_code,omitempty"`
	LatencyMs  int       `json:"latency_ms"`
	Success    bool      `json:"success"`
	Error      string    `json:"error"`
	CreatedAt  time.Time `json:"created_at"`
}

type InsertCheckInput struct {
	TargetID   string
	StatusCode *int
	LatencyMs  int
	Success    bool
	Error      string
}

type TargetStatus struct {
	Target
	LastCheck *Check `json:"last_check,omitempty"`
}
