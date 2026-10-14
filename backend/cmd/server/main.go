package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sam-ctd/job-tracker-agent/internal/agent"
	"github.com/sam-ctd/job-tracker-agent/internal/api"
	"github.com/sam-ctd/job-tracker-agent/internal/config"
	"github.com/sam-ctd/job-tracker-agent/internal/db"
	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	st := store.New(pool)
	agentOpts := agent.Options{
		LLM: agent.LLMConfig{
			APIKey:  cfg.LLMAPIKey,
			BaseURL: cfg.LLMBaseURL,
			Model:   cfg.LLMModel,
		},
	}
	if cfg.LLMAPIKey != "" {
		log.Printf("agent brain: llm (%s)", cfg.LLMModel)
	} else {
		log.Printf("agent brain: rules (set LLM_API_KEY to enable LLM)")
	}

	srv := api.New(st, pool, cfg.APIKey, cfg.CORSOrigins, agentOpts)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
