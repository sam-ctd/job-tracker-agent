package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sam-ctd/job-tracker-agent/internal/db"
)

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			http.Error(w, `{"status":"not_ready"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	addr := ":8088"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
