package http_pg_test

import (
    "context"
    "log"
    "os"
    "testing"
    "time"

    "github.com/fightingBald/GoTuto/internal/testutil"
)

var (
    pgDSN  string
    pgTemp bool
    pgCleanup func()
)

func TestMain(m *testing.M) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    dsn, isTemp, cleanup, err := testutil.DSNFromEnvOrDockerMain(ctx)
    if err != nil {
        log.Fatalf("pg fixture: %v", err)
    }
    pgDSN, pgTemp, pgCleanup = dsn, isTemp, cleanup

    code := m.Run()

    if pgCleanup != nil { pgCleanup() }
    os.Exit(code)
}

