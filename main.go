package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	cal "github.com/tom-fitz/score-keep-api/calendar"
	"github.com/tom-fitz/score-keep-api/imports"
	"github.com/tom-fitz/score-keep-api/league"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := &application{
		config: cfg,
		logger: logger,
	}

	jsonFile, err := os.Open("gcp-creds.json")
	if err != nil {
		log.Fatalf("Unable to get credentials file: %v", err)
	}
	// Read the contents of the file into a byte slice
	gcpCredsBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}
	gcpCreds, err := google.CredentialsFromJSON(ctx, gcpCredsBytes, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	gcpSvc, err := calendar.NewService(ctx, option.WithCredentials(gcpCreds))
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	router := mux.NewRouter()

	calendarHandler := cal.NewHandler(ctx, app.logger, 1, db, gcpSvc)
	calendarHandler.RegisterRoutes(router)

	importHandler := imports.NewHandler(ctx, app.logger, 1, db)
	importHandler.RegisterRoutes(router)

	leagueHandler := league.NewHandler(ctx, app.logger, 1, db)
	leagueHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	app.logger.Printf("Starting %s server at %s", cfg.env, addr)

	err = srv.ListenAndServe()
	app.logger.Fatal(err)
}
