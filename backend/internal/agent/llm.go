package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMConfig for OpenAI-compatible Chat Completions.
type LLMConfig struct {
	APIKey  string
	BaseURL string // e.g. https://api.openai.com/v1
	Model   string
}

// LLMBrain asks the model for a JSON tool decision.
type LLMBrain struct {
	cfg    LLMConfig
	client *http.Client
	fallback Brain
}

func NewLLMBrain(cfg LLMConfig, fallback Brain) *LLMBrain {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-4o-mini"
	}
	if fallback == nil {
		fallback = RulesBrain{}
	}
	return &LLMBrain{
		cfg:      cfg,
		client:   &http.Client{Timeout: 45 * time.Second},
		fallback: fallback,
	}
}

func (b *LLMBrain) Decide(ctx context.Context, userMessage string) (Decision, error) {
	dec, err := b.decideLLM(ctx, userMessage)
	if err != nil {
		// Network / parse failure → rules so chat keeps working.
		fb, fbErr := b.fallback.Decide(ctx, userMessage)
		if fbErr != nil {
			return Decision{}, fmt.Errorf("llm: %w; rules: %v", err, fbErr)
		}
		if fb.Message == "" {
			fb.Message = "(fallback rules; LLM unavailable)"
		} else {
			fb.Message = "(fallback rules) " + fb.Message
		}
		return fb, nil
	}
	return dec, nil
}

func (b *LLMBrain) decideLLM(ctx context.Context, userMessage string) (Decision, error) {
	body := map[string]any{
		"model": b.cfg.Model,
		"messages": []map[string]string{
			{"role": "system", "content": SystemPrompt()},
			{"role": "user", "content": userMessage},
		},
		"temperature": 0,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return Decision{}, err
	}

	url := strings.TrimRight(b.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return Decision{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.cfg.APIKey)

	res, err := b.client.Do(req)
	if err != nil {
		return Decision{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return Decision{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Decision{}, fmt.Errorf("llm http %d: %s", res.StatusCode, truncate(string(respBody), 300))
	}

	content, err := extractAssistantContent(respBody)
	if err != nil {
		return Decision{}, err
	}
	parsed, err := parseLLMJSON(content)
	if err != nil {
		return Decision{}, err
	}
	parsed.Source = "llm"
	return parsed, nil
}

func extractAssistantContent(respBody []byte) (string, error) {
	var envelope struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return "", err
	}
	if len(envelope.Choices) == 0 {
		return "", fmt.Errorf("llm: empty choices")
	}
	content := strings.TrimSpace(envelope.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("llm: empty content")
	}
	return content, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
