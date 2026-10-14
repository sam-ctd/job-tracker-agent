package agent

import (
	"fmt"
	"strings"
)

// SystemPrompt is sent to the LLM (tools stay the same as Catalog).
func SystemPrompt() string {
	var b strings.Builder
	b.WriteString(`Eres el asistente del Job Tracker. Ayudas a gestionar candidaturas.
Usa SOLO las tools listadas. No inventes IDs ni datos.
Statuses válidos: idea, enviada, relance, entretien, offre, refus.
Channels (slugs): linkedin, wttj, indeed, hellowork, apec, jobteaser, lesjeudis, chooseyourboss, alumni_estia.
Responde en el idioma del usuario (ES/FR/EN).

IMPORTANTE: responde ÚNICAMENTE con un JSON (sin markdown), uno de estos formatos:
1) Llamar tool: {"tool":"<name>","args":{...},"message":""}
2) Solo chat (sin tool): {"tool":null,"args":{},"message":"texto al usuario"}

Tools:
`)
	for _, t := range Catalog() {
		b.WriteString(fmt.Sprintf("- %s: %s\n", t.Name, t.Description))
		if len(t.Params) == 0 {
			continue
		}
		b.WriteString("  params:\n")
		for k, hint := range t.Params {
			b.WriteString(fmt.Sprintf("    %s: %s\n", k, hint))
		}
	}
	return b.String()
}
