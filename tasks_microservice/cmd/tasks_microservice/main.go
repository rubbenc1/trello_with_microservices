package main

import (
    "database/sql"
    "fmt"
    "log/slog"
    "os"
    "strconv"
    "tasks_microservice/internal/config"
    "tasks_microservice/internal/database"
    "tasks_microservice/internal/handlers"
    "tasks_microservice/internal/rabbitmq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
    envLocal = "local"
    envDev   = "dev"
)

func main() {
    cfg := config.MustLoad()
    log := setupLogger(cfg.Env)
    log.Info("starting a microservice", slog.String("env", cfg.Env))

    // Connect to the database
    db, err := database.Open(&cfg.Postgres)
    if err != nil {
        log.Error("Failed to connect to database", slog.String("error", err.Error()))
        os.Exit(1)
    }
    defer db.Close()

    if err := runMigrations(db.Conn, cfg.Postgres); err != nil {
        log.Error("failed to run migrations", slog.String("error", err.Error()))
        return
    }

    // Connect to RabbitMQ
    rabbit, err := rabbitmq.ConnectRabbitMQ("amqp://" + cfg.RabbitMQ.User + ":" + cfg.RabbitMQ.Password + "@" + cfg.RabbitMQ.Host + ":" + strconv.Itoa(cfg.RabbitMQ.RabbitMQPort) + "/")
    if err != nil {
        log.Error("Failed to connect to RabbitMQ", slog.String("error", err.Error()))
        os.Exit(1)
    }
    defer rabbit.Close()

    // Set up handlers with the database connection
    handler := handlers.NewHandlers(db, rabbit)

    // Set up RabbitMQ consumers
    rabbit.SetupRabbitMQConsumers(handler)

    select {}
}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger
    switch env {
    case envLocal:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envDev:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    default:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
        log.Warn("unsupported environment, defaulting to text logger", slog.String("env", env))
    }

    return log
}

func runMigrations(db *sql.DB, cfg config.Postgres) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("failed to create migration driver: %w", err)
    }

    migrationsPath := "file://../../migrations"
    m, err := migrate.NewWithDatabaseInstance(
        migrationsPath,
        cfg.DBName, driver)
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to apply migrations: %w", err)
    }

    return nil
}
