package agent

import (
	"fmt"
	"strings"
)

// SystemPrompt is what 3c will send to the LLM (tools stay the same).
func SystemPrompt() string {
	var b strings.Builder
	b.WriteString(`Eres el asistente del Job Tracker. Ayudas a gestionar candidaturas.
Usa SOLO las tools listadas. No inventes IDs ni datos.
Statuses válidos: idea, enviada, relance, entretien, offre, refus.
Channels (slugs): linkedin, wttj, indeed, hellowork, apec, jobteaser, lesjeudis, chooseyourboss, alumni_estia.
Responde en el idioma del usuario (ES/FR/EN).

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
