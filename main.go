package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-postgres-boilerplate/app"
	"go-postgres-boilerplate/dao"
)

func main() {
	cfg := dao.NewConfig()
	database, err := dao.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	personDAO, err := dao.NewPersonDAO(database)
	if err != nil {
		log.Fatalf("Failed to create PersonDAO: %v", err)
	}

	application := app.NewApp(personDAO)
	router := application.SetupRouter()

	address := fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT"))
	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting App on %s", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
