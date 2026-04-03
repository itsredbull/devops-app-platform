package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dsn string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *PostgresStore) ApplyMigrations(ctx context.Context, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		path := filepath.Join(migrationsDir, name)
		query, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read migration %s: %w", name, readErr)
		}
		if _, execErr := s.pool.Exec(ctx, string(query)); execErr != nil {
			return fmt.Errorf("exec migration %s: %w", name, execErr)
		}
	}

	return nil
}

func (s *PostgresStore) CreateTarget(ctx context.Context, in CreateTargetInput) (Target, error) {
	query := `
		INSERT INTO targets (url, check_interval_seconds, timeout_seconds, enabled)
		VALUES ($1, $2, $3, $4)
		RETURNING id, url, check_interval_seconds, timeout_seconds, enabled, created_at, updated_at
	`

	var t Target
	err := s.pool.QueryRow(ctx, query, in.URL, in.CheckIntervalSeconds, in.TimeoutSeconds, in.Enabled).Scan(
		&t.ID,
		&t.URL,
		&t.CheckIntervalSeconds,
		&t.TimeoutSeconds,
		&t.Enabled,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return Target{}, err
	}

	return t, nil
}

func (s *PostgresStore) GetTarget(ctx context.Context, id string) (Target, error) {
	query := `
		SELECT id, url, check_interval_seconds, timeout_seconds, enabled, created_at, updated_at
		FROM targets
		WHERE id = $1
	`

	var t Target
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.URL,
		&t.CheckIntervalSeconds,
		&t.TimeoutSeconds,
		&t.Enabled,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Target{}, ErrNotFound
		}
		return Target{}, err
	}

	return t, nil
}

func (s *PostgresStore) ListTargets(ctx context.Context) ([]Target, error) {
	query := `
		SELECT id, url, check_interval_seconds, timeout_seconds, enabled, created_at, updated_at
		FROM targets
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Target, 0)
	for rows.Next() {
		var t Target
		err = rows.Scan(&t.ID, &t.URL, &t.CheckIntervalSeconds, &t.TimeoutSeconds, &t.Enabled, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, rows.Err()
}

func (s *PostgresStore) ListEnabledTargets(ctx context.Context) ([]Target, error) {
	query := `
		SELECT id, url, check_interval_seconds, timeout_seconds, enabled, created_at, updated_at
		FROM targets
		WHERE enabled = TRUE
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Target
	for rows.Next() {
		var t Target
		err = rows.Scan(&t.ID, &t.URL, &t.CheckIntervalSeconds, &t.TimeoutSeconds, &t.Enabled, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *PostgresStore) DeleteTarget(ctx context.Context, id string) (bool, error) {
	cmd, err := s.pool.Exec(ctx, `DELETE FROM targets WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}

func (s *PostgresStore) TruncateAll(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `TRUNCATE TABLE checks, targets RESTART IDENTITY CASCADE`)
	return err
}

func (s *PostgresStore) InsertCheck(ctx context.Context, in InsertCheckInput) error {
	query := `
		INSERT INTO checks (target_id, status_code, latency_ms, success, error)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.pool.Exec(ctx, query, in.TargetID, in.StatusCode, in.LatencyMs, in.Success, in.Error)
	return err
}

func (s *PostgresStore) ListChecks(ctx context.Context, targetID string, limit int) ([]Check, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	query := `
		SELECT id, target_id, checked_at, status_code, latency_ms, success, COALESCE(error, ''), created_at
		FROM checks
		WHERE ($1 = '' OR target_id = $1)
		ORDER BY checked_at DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, query, targetID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Check
	for rows.Next() {
		var c Check
		err = rows.Scan(&c.ID, &c.TargetID, &c.CheckedAt, &c.StatusCode, &c.LatencyMs, &c.Success, &c.Error, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *PostgresStore) ListLatestStatus(ctx context.Context) ([]TargetStatus, error) {
	query := `
		SELECT
			t.id,
			t.url,
			t.check_interval_seconds,
			t.timeout_seconds,
			t.enabled,
			t.created_at,
			t.updated_at,
			c.id,
			c.target_id,
			c.checked_at,
			c.status_code,
			c.latency_ms,
			c.success,
			COALESCE(c.error, ''),
			c.created_at
		FROM targets t
		LEFT JOIN LATERAL (
			SELECT id, target_id, checked_at, status_code, latency_ms, success, error, created_at
			FROM checks
			WHERE target_id = t.id
			ORDER BY checked_at DESC
			LIMIT 1
		) c ON TRUE
		ORDER BY t.created_at DESC
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TargetStatus
	for rows.Next() {
		var ts TargetStatus
		var checkID *string
		var targetID *string
		var checkedAt *time.Time
		var statusCode *int
		var latencyMs *int
		var success *bool
		var errText *string
		var checkCreatedAt *time.Time

		err = rows.Scan(
			&ts.ID,
			&ts.URL,
			&ts.CheckIntervalSeconds,
			&ts.TimeoutSeconds,
			&ts.Enabled,
			&ts.CreatedAt,
			&ts.UpdatedAt,
			&checkID,
			&targetID,
			&checkedAt,
			&statusCode,
			&latencyMs,
			&success,
			&errText,
			&checkCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if checkID != nil && targetID != nil && checkedAt != nil && latencyMs != nil && success != nil && checkCreatedAt != nil {
			check := Check{
				ID:         *checkID,
				TargetID:   *targetID,
				CheckedAt:  *checkedAt,
				StatusCode: statusCode,
				LatencyMs:  *latencyMs,
				Success:    *success,
				CreatedAt:  *checkCreatedAt,
			}
			if errText != nil {
				check.Error = *errText
			}
			ts.LastCheck = &check
		}

		out = append(out, ts)
	}

	return out, rows.Err()
}
