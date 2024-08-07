package database

import (
    "database/sql"
    "fmt"
    "log"
    "tasks_microservice/internal/config"

    _ "github.com/lib/pq"
)

type Database struct {
    Conn *sql.DB
}

func Open(cfg *config.Postgres) (*Database, error) {
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }

    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &Database{Conn: db}, nil
}

func (db *Database) Close() {
    if err := db.Conn.Close(); err != nil {
        log.Printf("failed to close database connection: %v", err)
    }
}
