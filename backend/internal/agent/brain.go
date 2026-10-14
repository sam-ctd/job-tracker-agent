package agent

import "context"

// Decision is what a brain (rules or LLM) returns before tools run.
type Decision struct {
	Call    *Call  // nil → no tool, just Message
	Message string // optional text (LLM chat or hint)
	Source  string // "llm" | "rules"
}

// Brain chooses a tool call (or a plain reply).
type Brain interface {
	Decide(ctx context.Context, userMessage string) (Decision, error)
}

// RulesBrain is the keyword router from 3a/3b.
type RulesBrain struct{}

func (RulesBrain) Decide(_ context.Context, userMessage string) (Decision, error) {
	call, ok := decide(userMessage)
	if !ok {
		return Decision{
			Message: "No sé qué tool usar (modo rules). Prueba stats, lista, due, crea o patch.",
			Source:  "rules",
		}, nil
	}
	c := call
	return Decision{Call: &c, Source: "rules"}, nil
}
