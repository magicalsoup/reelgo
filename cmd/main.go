package main

import (
	"context"
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
	"github.com/magicalsoup/reelgo-backend/src/instagram"
)

func NewServer() http.Handler {
	r := mux.NewRouter()
	instagram.AddRoutes(r)
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

	server := NewServer()
	httpServer := &http.Server{
		Addr: net.JoinHostPort("localhost", "8080"),
		Handler: server,
	}

	go func() {
		log.Printf("listening on server %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
