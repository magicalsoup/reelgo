package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/magicalsoup/reelgo/src/instagram"

	_ "github.com/lib/pq"
)

func NewServer(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	instagram.AddRoutes(r, db)
	var handler http.Handler = r
	return handler
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL_NEON"))
	
	if err != nil {
		fmt.Println("unable to open database")
		return err
	}

	server := NewServer(db)

	httpServer := &http.Server{
		Addr: net.JoinHostPort(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		Handler: server,
	}

	go func() {
		log.Printf("listening on server %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1) // add go routine counter
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	
	wg.Wait() // wait for go routines to finish
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
