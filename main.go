// Package main is the postgresql lesson `l1_acid_and_mvcc` homework scaffold for Vibe Learn.
//
// Задача: MVCC-демо: 1 писатель + 3 долгих читателя, замер n_dead_tup и отсутствия read-locks.
// Реализуй функции ниже — сигнатуры и тестовая поверхность фиксированы;
// CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// Драйвер: github.com/jackc/pgx/v5 (+ pgxpool). DATABASE_URL — DSN из env.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Latencies — собранные перцентили для бенчмарка запроса.
type Latencies struct{ P50, P95, P99 time.Duration }

// StandbyInfo — строка из pg_stat_replication для выбора кандидата на promote.
type StandbyInfo struct {
	ClientAddr string
	ReplayLSN  string
	State      string
}

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DatabaseURL — DSN PostgreSQL. Дефолт совпадает с docker-compose.yml.
func DatabaseURL() string {
	return envOr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
}

// Connect открывает пул pgx из DATABASE_URL.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, DatabaseURL())
}

// ----- TODO #1: DeadTuples -----
//
// SELECT n_dead_tup FROM pg_stat_user_tables WHERE relname=$1
func DeadTuples(ctx context.Context, pool *pgxpool.Pool, table string) (int64, error) {
	// TODO: implement
	panic("DeadTuples: not implemented")
}

// ----- TODO #2: RunWriterReaders -----
//
// 1 goroutine UPDATE в цикле + 3 goroutine SELECT в долгих транзакциях (Repeatable Read)
func RunWriterReaders(ctx context.Context, pool *pgxpool.Pool, rowID int64, iters int) error {
	// TODO: implement
	panic("RunWriterReaders: not implemented")
}

// ----- TODO #3: BlockingLocks -----
//
// SELECT count(*) FROM pg_locks WHERE NOT granted — должно остаться 0
func BlockingLocks(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	// TODO: implement
	panic("BlockingLocks: not implemented")
}

// _refs keeps imports live while the TODO bodies are unimplemented stubs.
// Удали эту функцию, когда реализуешь TODO выше.
var _refs = []any{
	Latencies{},
	StandbyInfo{},
	time.Second,
}

// ----- main entry -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — postgresql lesson %s scaffold up", "l1_acid_and_mvcc")
	log.Printf("DATABASE_URL: %s", DatabaseURL())
	log.Printf("Реализуй TODO-функции, затем `go test ./...`. README.md содержит задачу.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown so `go run .` is interactive — Ctrl-C exits cleanly.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutdown signal received")
		cancel()
	}()
	<-ctx.Done()
}
