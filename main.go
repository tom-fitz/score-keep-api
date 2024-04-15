package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/tom-fitz/score-keep-api/imports"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port   int
	env    string
	dbConn string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "environment (dev|prod)")
	flag.StringVar(&cfg.dbConn, "db-conn", "postgres://admin:admin@localhost:5432/score_keep_db?sslmode=disable", "PostgreSQL database DSN")
	flag.Parse()

	db, dbErr := sql.Open("postgres", cfg.dbConn)
	if dbErr != nil {
		log.Fatal("could not connect to database:", dbErr)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	importHandler := imports.NewHandler(app.logger, 1, db)

	http.Handle("/v1/import/", addCorsHeaders(importHandler))

	srv := &http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("Starting %s server at %s", cfg.env, addr)

	err := srv.ListenAndServe()
	app.logger.Fatal(err)
}

func addCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("%s, %s, %s, %s, %s", http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
