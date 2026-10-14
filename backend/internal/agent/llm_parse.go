package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type llmJSON struct {
	Tool    *string        `json:"tool"`
	Args    map[string]any `json:"args"`
	Message string         `json:"message"`
}

func parseLLMJSON(content string) (Decision, error) {
	content = stripCodeFence(strings.TrimSpace(content))
	if !strings.HasPrefix(content, "{") {
		if i := strings.Index(content, "{"); i >= 0 {
			if j := strings.LastIndex(content, "}"); j > i {
				content = content[i : j+1]
			}
		}
	}

	var raw llmJSON
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return Decision{}, fmt.Errorf("parse llm json: %w", err)
	}

	if raw.Tool == nil || strings.TrimSpace(*raw.Tool) == "" || strings.EqualFold(*raw.Tool, "null") {
		msg := strings.TrimSpace(raw.Message)
		if msg == "" {
			msg = "¿En qué te ayudo con tus candidaturas?"
		}
		return Decision{Message: msg, Source: "llm"}, nil
	}

	name := strings.TrimSpace(*raw.Tool)
	if !isKnownTool(name) {
		return Decision{}, fmt.Errorf("unknown tool %q", name)
	}
	call := Call{Name: name, Args: stringifyArgs(raw.Args)}
	return Decision{Call: &call, Message: strings.TrimSpace(raw.Message), Source: "llm"}, nil
}

func stringifyArgs(in map[string]any) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		if v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			out[k] = t
		case float64:
			// JSON numbers
			if t == float64(int64(t)) {
				out[k] = fmt.Sprintf("%d", int64(t))
			} else {
				out[k] = fmt.Sprintf("%v", t)
			}
		case bool:
			out[k] = fmt.Sprintf("%t", t)
		default:
			out[k] = fmt.Sprintf("%v", t)
		}
	}
	return out
}

func isKnownTool(name string) bool {
	for _, t := range Catalog() {
		if t.Name == name {
			return true
		}
	}
	return false
}

func stripCodeFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSpace(s)
	if strings.HasPrefix(strings.ToLower(s), "json") {
		s = strings.TrimSpace(s[4:])
	}
	if i := strings.LastIndex(s, "```"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
