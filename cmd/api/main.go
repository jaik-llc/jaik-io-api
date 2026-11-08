package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jaik-llc/jaik-io-api/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	port, err := strconv.Atoi(envOr("PORT", "8080"))
	if err != nil {
		logger.Fatalf("invalid PORT: %v", err)
	}
	// The connection string comes from the environment (a Kubernetes Secret in production),
	// never from source code.
	cfg := config{
		port: port,
		env:  envOr("ENV", "dev"),
		dsn:  os.Getenv("READINGLIST_DB_DSN"),
	}

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	// sql.Open does not connect. Ping once so a bad DSN shows up in the logs at startup, but keep
	// serving: /healthcheck/db reports the database state and the pool reconnects on its own.
	if err := db.Ping(); err != nil {
		logger.Printf("database not reachable at startup: %v", err)
	} else {
		logger.Printf("database connection pool established")
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
