package agent

import (
	"context"
	"strings"

	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

// Agent runs a rule-based brain that picks tools and calls the store.
// 3c: replace decide() with an LLM; Catalog/SystemPrompt/runTool stay.
type Agent struct {
	store *store.Store
}

func New(st *store.Store) *Agent {
	return &Agent{store: st}
}

// Reply is what /api/chat returns.
type Reply struct {
	Message  string            `json:"message"`
	Tool     string            `json:"tool,omitempty"`
	Args     map[string]string `json:"args,omitempty"`
	Result   any               `json:"result,omitempty"`
	HelpHint string            `json:"help_hint,omitempty"`
}

func (a *Agent) Chat(ctx context.Context, userMessage string) (Reply, error) {
	msg := strings.TrimSpace(userMessage)
	if msg == "" {
		return Reply{
			Message:  "Escribe un mensaje. Puedo listar, stats, due, crear y actualizar candidaturas.",
			HelpHint: helpText(),
		}, nil
	}

	call, ok := decide(msg)
	if !ok {
		return Reply{
			Message:  "No sé qué tool usar todavía (brain fake, sin LLM).",
			HelpHint: helpText(),
		}, nil
	}

	if call.Name == toolCreate && (call.Args["company"] == "" || call.Args["title"] == "") {
		return Reply{
			Message:  "Para crear hace falta company y title.",
			Tool:     toolCreate,
			Args:     call.Args,
			HelpHint: `Ejemplo: crea Acme | Backend Go   o   crear company=Acme title="Dev Go"`,
		}, nil
	}
	if call.Name == toolPatch && call.Args["id"] == "" {
		return Reply{
			Message:  "Para actualizar hace falta el id (UUID).",
			Tool:     toolPatch,
			Args:     call.Args,
			HelpHint: `Ejemplo: patch id=<uuid> status=enviada`,
		}, nil
	}
	if call.Name == toolPatch && call.Args["status"] == "" && call.Args["company"] == "" &&
		call.Args["title"] == "" && call.Args["notes"] == "" && call.Args["channel"] == "" {
		return Reply{
			Message:  "Di qué campo cambiar (p.ej. status=relance).",
			Tool:     toolPatch,
			Args:     call.Args,
			HelpHint: `Ejemplo: actualiza <uuid> status=entretien`,
		}, nil
	}

	result, text, err := runTool(ctx, a.store, call)
	if err != nil {
		return Reply{}, err
	}
	return Reply{
		Message: text,
		Tool:    call.Name,
		Args:    call.Args,
		Result:  result,
	}, nil
}

func helpText() string {
	return strings.Join([]string{
		`stats → "muéstrame las stats"`,
		`lista → "lista las candidaturas"`,
		`due → "qué hay due / relances"`,
		`crea → "crea Acme | Backend Go"`,
		`patch → "patch id=<uuid> status=enviada"`,
	}, " · ")
}
