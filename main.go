package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tom-fitz/score-keep-api/imports"
	"github.com/tom-fitz/score-keep-api/league"

	_ "github.com/lib/pq"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	dbConnStr := os.Getenv("DB_CONN")
	if dbConnStr == "" {
		dbConnStr = "postgres://admin:admin@localhost:5432/score_keep_db?sslmode=disable"
	}

	db, dbErr := sql.Open("postgres", dbConnStr)
	if dbErr != nil {
		log.Fatalf("could not connect to database: %v", dbErr)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//gcpSvc, err := cal.NewCalendarService(ctx, "gcp-creds.json")
	//if err != nil {
	//	log.Fatalf("could not create GCP Calendar service: %v", err)
	//}

	addr := fmt.Sprintf(":%s", port)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	//calendarHandler := cal.NewService(ctx, log, 1, db, gcpSvc)
	//calendarHandler.RegisterRoutes(router)

	importHandler := imports.NewService(ctx, log, 1, db)
	importHandler.RegisterRoutes(router)

	leagueHandler := league.NewHandler(ctx, log, 1, db)
	leagueHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	go func() {
		log.Infof("Starting server at %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Infof("Received signal: %v, shutting down gracefully...", sig)

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Info("Graceful shutdown completed")
}
