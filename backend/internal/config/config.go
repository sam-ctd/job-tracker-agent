package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	Port        string
	APIKey      string
	CORSOrigins []string

	// Optional LLM (OpenAI-compatible). Empty API key → rules brain only.
	LLMAPIKey  string
	LLMBaseURL string
	LLMModel   string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		APIKey:      os.Getenv("APP_API_KEY"),
		LLMAPIKey:   os.Getenv("LLM_API_KEY"),
		LLMBaseURL:  os.Getenv("LLM_BASE_URL"),
		LLMModel:    os.Getenv("LLM_MODEL"),
	}
	if cfg.Port == "" {
		cfg.Port = "8088"
	}
	if cfg.LLMBaseURL == "" {
		cfg.LLMBaseURL = "https://api.openai.com/v1"
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = "gpt-4o-mini"
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.APIKey == "" {
		return cfg, fmt.Errorf("APP_API_KEY is required")
	}
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		raw = "http://localhost:5173,http://127.0.0.1:5173,http://localhost:3847"
	}
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			cfg.CORSOrigins = append(cfg.CORSOrigins, o)
		}
	}
	return cfg, nil
}
