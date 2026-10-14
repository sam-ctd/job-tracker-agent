package agent

import "strings"

const (
	toolGetStats          = "get_stats"
	toolListApplications  = "list_applications"
)

// decideTool is the fake brain: keyword matching only (FR/ES/EN).
func decideTool(message string) string {
	m := strings.ToLower(message)

	switch {
	case containsAny(m, "stat", "stats", "estad", "résumé", "resume", "récap", "recap", "cuánt", "cuant", "combien", "total"):
		return toolGetStats
	case containsAny(m, "list", "lista", "liste", "candidatur", "application", "muestra", "montre", "show"):
		return toolListApplications
	default:
		return ""
	}
}

func containsAny(s string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}
