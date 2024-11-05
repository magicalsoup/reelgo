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
	"github.com/magicalsoup/reelgo/src/auth"
	"github.com/magicalsoup/reelgo/src/instagram"
	"github.com/magicalsoup/reelgo/src/users"

	_ "github.com/lib/pq"
)
// CORS Middleware to handle preflight requests and set CORS headers
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

        // Handle preflight OPTIONS request
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent) // Return 204 No Content for OPTIONS requests
            return
        }

        next.ServeHTTP(w, r) // Pass the request to the next handler
    })
}

func NewServer(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	instagram.AddRoutes(r, db)
	users.AddRoutes(r, db)
	auth.AddRoutes(r, db)
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
