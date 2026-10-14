package agent

import (
	"context"
	"fmt"

	"github.com/sam-ctd/job-tracker-agent/internal/models"
	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

const (
	toolGetStats         = "get_stats"
	toolListApplications = "list_applications"
	toolListDue          = "list_due"
	toolCreate           = "create_application"
	toolPatch            = "patch_application"
)

// ToolSpec describes a tool for the fake brain and (later) the LLM.
type ToolSpec struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Params      map[string]string `json:"params"` // name → short type/hint
}

// Call is what the brain (fake or LLM) decides to run.
type Call struct {
	Name string
	Args map[string]string
}

// Catalog is the stable tool list (3c will send this to the model).
func Catalog() []ToolSpec {
	return []ToolSpec{
		{
			Name:        toolGetStats,
			Description: "Resumen: total y conteos por status/channel.",
			Params:      map[string]string{},
		},
		{
			Name:        toolListApplications,
			Description: "Lista candidaturas. Filtros opcionales status y company.",
			Params: map[string]string{
				"status":  "string? one of idea|enviada|relance|entretien|offre|refus",
				"company": "string? substring ILIKE",
			},
		},
		{
			Name:        toolListDue,
			Description: "Candidaturas con next_follow_up vencido (excluye offre/refus).",
			Params:      map[string]string{},
		},
		{
			Name:        toolCreate,
			Description: "Crea una candidatura. company y title son obligatorios.",
			Params: map[string]string{
				"company":  "string required",
				"title":    "string required",
				"url":      "string?",
				"channel":  "string? default linkedin",
				"status":   "string? default idea",
				"notes":    "string?",
				"contact":  "string?",
				"location": "string?",
				"remote":   "string?",
				"stack":    "string?",
			},
		},
		{
			Name:        toolPatch,
			Description: "Actualiza campos de una candidatura por id (UUID).",
			Params: map[string]string{
				"id":     "uuid required",
				"status": "string?",
				"company": "string?",
				"title":  "string?",
				"notes":  "string?",
				"channel": "string?",
			},
		},
	}
}

func runTool(ctx context.Context, st *store.Store, call Call) (any, string, error) {
	switch call.Name {
	case toolGetStats:
		stats, err := st.Stats(ctx)
		if err != nil {
			return nil, "", err
		}
		return stats, fmt.Sprintf("Tienes %d candidatura(s) en total.", stats.Total), nil

	case toolListApplications:
		apps, err := st.List(ctx, call.Args["status"], call.Args["company"])
		if err != nil {
			return nil, "", err
		}
		return apps, fmt.Sprintf("Listado: %d candidatura(s).", len(apps)), nil

	case toolListDue:
		apps, err := st.Due(ctx)
		if err != nil {
			return nil, "", err
		}
		return apps, fmt.Sprintf("Pendientes de follow-up: %d.", len(apps)), nil

	case toolCreate:
		in := models.ApplicationCreate{
			Company:  call.Args["company"],
			Title:    call.Args["title"],
			URL:      call.Args["url"],
			Channel:  call.Args["channel"],
			Status:   call.Args["status"],
			Notes:    call.Args["notes"],
			Contact:  call.Args["contact"],
			Location: call.Args["location"],
			Remote:   call.Args["remote"],
			Stack:    call.Args["stack"],
		}
		app, err := st.Create(ctx, in)
		if err != nil {
			return nil, "", err
		}
		return app, fmt.Sprintf("Creada: %s @ %s (%s).", app.Title, app.Company, app.ID), nil

	case toolPatch:
		id := call.Args["id"]
		patch := models.ApplicationPatch{}
		if v, ok := call.Args["status"]; ok && v != "" {
			patch.Status = &v
		}
		if v, ok := call.Args["company"]; ok && v != "" {
			patch.Company = &v
		}
		if v, ok := call.Args["title"]; ok && v != "" {
			patch.Title = &v
		}
		if v, ok := call.Args["notes"]; ok && v != "" {
			patch.Notes = &v
		}
		if v, ok := call.Args["channel"]; ok && v != "" {
			patch.Channel = &v
		}
		app, err := st.Patch(ctx, id, patch)
		if err != nil {
			return nil, "", err
		}
		return app, fmt.Sprintf("Actualizada %s → status=%s.", app.ID, app.Status), nil

	default:
		return nil, "", fmt.Errorf("unknown tool %q", call.Name)
	}
}
