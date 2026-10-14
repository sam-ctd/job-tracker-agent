package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

// Agent runs a tiny rule-based "brain" that picks tools and calls the store.
// Later (3c) a real LLM replaces decideTool; the tools stay the same.
type Agent struct {
	store *store.Store
}

func New(st *store.Store) *Agent {
	return &Agent{store: st}
}

// Reply is what /api/chat returns.
type Reply struct {
	Message  string `json:"message"`
	Tool     string `json:"tool,omitempty"`
	Result   any    `json:"result,omitempty"`
	HelpHint string `json:"help_hint,omitempty"`
}

func (a *Agent) Chat(ctx context.Context, userMessage string) (Reply, error) {
	msg := strings.TrimSpace(userMessage)
	if msg == "" {
		return Reply{
			Message:  "Escribe un mensaje. Por ahora entiendo stats y listar candidaturas.",
			HelpHint: helpText(),
		}, nil
	}

	tool := decideTool(msg)
	switch tool {
	case toolGetStats:
		stats, err := a.store.Stats(ctx)
		if err != nil {
			return Reply{}, err
		}
		return Reply{
			Message: fmt.Sprintf("Tienes %d candidatura(s) en total.", stats.Total),
			Tool:    toolGetStats,
			Result:  stats,
		}, nil

	case toolListApplications:
		apps, err := a.store.List(ctx, "", "")
		if err != nil {
			return Reply{}, err
		}
		return Reply{
			Message: fmt.Sprintf("Listado: %d candidatura(s).", len(apps)),
			Tool:    toolListApplications,
			Result:  apps,
		}, nil

	default:
		return Reply{
			Message:  "No sé qué tool usar todavía (brain fake, sin LLM). Prueba con stats o lista.",
			HelpHint: helpText(),
		}, nil
	}
}

func helpText() string {
	return "Ejemplos: \"muéstrame las stats\", \"lista las candidaturas\""
}
