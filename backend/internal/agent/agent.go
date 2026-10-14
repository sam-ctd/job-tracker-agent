package agent

import (
	"context"
	"strings"

	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

// Agent: brain elige tool → runTool ejecuta sobre el store.
type Agent struct {
	store *store.Store
	brain Brain
}

type Options struct {
	LLM LLMConfig // si APIKey vacío → solo RulesBrain
}

func New(st *store.Store, opts Options) *Agent {
	var brain Brain = RulesBrain{}
	if strings.TrimSpace(opts.LLM.APIKey) != "" {
		brain = NewLLMBrain(opts.LLM, RulesBrain{})
	}
	return &Agent{store: st, brain: brain}
}

// Reply is what /api/chat returns.
type Reply struct {
	Message  string            `json:"message"`
	Tool     string            `json:"tool,omitempty"`
	Args     map[string]string `json:"args,omitempty"`
	Result   any               `json:"result,omitempty"`
	HelpHint string            `json:"help_hint,omitempty"`
	Brain    string            `json:"brain,omitempty"` // "llm" | "rules"
}

func (a *Agent) Chat(ctx context.Context, userMessage string) (Reply, error) {
	msg := strings.TrimSpace(userMessage)
	if msg == "" {
		return Reply{
			Message:  "Escribe un mensaje. Puedo listar, stats, due, crear y actualizar candidaturas.",
			HelpHint: helpText(),
			Brain:    "rules",
		}, nil
	}

	dec, err := a.brain.Decide(ctx, msg)
	if err != nil {
		return Reply{}, err
	}

	if dec.Call == nil {
		out := Reply{Message: dec.Message, Brain: dec.Source, HelpHint: helpText()}
		if out.Message == "" {
			out.Message = "¿En qué te ayudo?"
		}
		return out, nil
	}

	call := *dec.Call
	if call.Name == toolCreate && (call.Args["company"] == "" || call.Args["title"] == "") {
		return Reply{
			Message:  "Para crear hace falta company y title.",
			Tool:     toolCreate,
			Args:     call.Args,
			HelpHint: `Ejemplo: crea Acme | Backend Go`,
			Brain:    dec.Source,
		}, nil
	}
	if call.Name == toolPatch && call.Args["id"] == "" {
		return Reply{
			Message:  "Para actualizar hace falta el id (UUID).",
			Tool:     toolPatch,
			Args:     call.Args,
			HelpHint: `Ejemplo: patch id=<uuid> status=enviada`,
			Brain:    dec.Source,
		}, nil
	}
	if call.Name == toolPatch && call.Args["status"] == "" && call.Args["company"] == "" &&
		call.Args["title"] == "" && call.Args["notes"] == "" && call.Args["channel"] == "" {
		return Reply{
			Message:  "Di qué campo cambiar (p.ej. status=relance).",
			Tool:     toolPatch,
			Args:     call.Args,
			HelpHint: `Ejemplo: actualiza <uuid> status=entretien`,
			Brain:    dec.Source,
		}, nil
	}

	result, text, err := runTool(ctx, a.store, call)
	if err != nil {
		return Reply{}, err
	}
	if dec.Message != "" {
		text = dec.Message + " " + text
	}
	return Reply{
		Message: text,
		Tool:    call.Name,
		Args:    call.Args,
		Result:  result,
		Brain:   dec.Source,
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
