package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"flag"
	"log"
	"os"
	"time"

	"github.com/jaik-llc/jaik-io-api/internal/data"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	Port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models

}


func main() {
	var cfg config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", "postgresql://readinglist:pa$$w0rd@k8s-readingl-postgres-1fb4f7366f-bb453d3c74b6f6ff.elb.us-east-2.amazonaws.com:5432/readinglist?sslmode=require", "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.Models{
		Book: data.BookModel{DB: db},
		},
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}



