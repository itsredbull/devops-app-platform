package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/api"
	"github.com/itsredbull/devops-app-platform/app/internal/checker"
	"github.com/itsredbull/devops-app-platform/app/internal/scheduler"
	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

func main() {
	port := getenv("APP_PORT", "8080")
	dsn := dbDSNFromEnv()
	migrationsDir := getenv("MIGRATIONS_DIR", "./migrations")
	autoMigrate := getenv("AUTO_MIGRATE", "true") == "true"
	tickSeconds := getenvInt("SCHEDULER_TICK_SECONDS", 5)
	maxAttempts := getenvInt("CHECK_MAX_ATTEMPTS", 3)
	retryBackoffMs := getenvInt("CHECK_RETRY_BACKOFF_MS", 200)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := store.NewPostgresStore(ctx, dsn)
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	defer st.Close()

	if autoMigrate {
		if err = st.ApplyMigrations(ctx, migrationsDir); err != nil {
			log.Fatalf("apply migrations failed: %v", err)
		}
	}

	apiServer := api.NewServer(st)
	mux := http.NewServeMux()
	apiServer.RegisterRoutes(mux)
	root := api.WithMetrics(mux)

	worker := scheduler.NewWorker(
		st,
		checker.NewService(maxAttempts, time.Duration(retryBackoffMs)*time.Millisecond),
		time.Duration(tickSeconds)*time.Second,
	)
	go worker.Start(ctx)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("uptime-api listening on :%s", port)
		if serveErr := srv.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", serveErr)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("shutdown complete")
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func dbDSNFromEnv() string {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "uptime")
	password := getenv("DB_PASSWORD", "uptime")
	name := getenv("DB_NAME", "uptime")
	sslmode := getenv("DB_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}
