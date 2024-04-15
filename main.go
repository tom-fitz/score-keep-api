package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tom-fitz/score-keep-api/imports"
	"github.com/tom-fitz/score-keep-api/league"
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

	router := mux.NewRouter()

	importHandler := imports.NewHandler(app.logger, 1, db)
	importHandler.RegisterRoutes(router)

	leagueHandler := league.NewHandler(app.logger, 1, db)
	leagueHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	app.logger.Printf("Starting %s server at %s", cfg.env, addr)

	err := srv.ListenAndServe()
	app.logger.Fatal(err)
}
