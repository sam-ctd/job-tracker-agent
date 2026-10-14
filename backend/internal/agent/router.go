package agent

import (
	"regexp"
	"strings"

	"github.com/sam-ctd/job-tracker-agent/internal/models"
)

var uuidRe = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// decide is the fake brain: keywords + simple arg parsing (FR/ES/EN).
func decide(message string) (Call, bool) {
	m := strings.TrimSpace(message)
	lower := strings.ToLower(m)

	switch {
	case containsAny(lower, "crea", "crear", "ajoute", "ajouter", "add ", "nueva", "nouveau", "new application"):
		return callCreate(m), true
	case containsAny(lower, "patch", "actualiza", "actualizar", "cambia", "cambiar", "update", "passe ", "pasar "):
		return callPatch(m, lower), true
	case containsAny(lower, "due", "follow", "relance", "relancer", "seguimiento", "à relancer", "a relancer", "pendiente"):
		return Call{Name: toolListDue, Args: map[string]string{}}, true
	case containsAny(lower, "stat", "stats", "estad", "récap", "recap", "cuánt", "cuant", "combien", "total"):
		return Call{Name: toolGetStats, Args: map[string]string{}}, true
	case containsAny(lower, "list", "lista", "liste", "muestra", "montre", "show", "candidatur", "application"):
		return Call{Name: toolListApplications, Args: parseFilters(m)}, true
	default:
		return Call{}, false
	}
}

func callCreate(message string) Call {
	args := parseKeyValues(message)
	// Fallback: "crea Acme | Backend Go" or "crear Acme / Backend"
	if args["company"] == "" || args["title"] == "" {
		rest := stripCreateVerb(message)
		if parts := splitCompanyTitle(rest); len(parts) == 2 {
			if args["company"] == "" {
				args["company"] = parts[0]
			}
			if args["title"] == "" {
				args["title"] = parts[1]
			}
		}
	}
	return Call{Name: toolCreate, Args: args}
}

func callPatch(message, lower string) Call {
	args := parseKeyValues(message)
	if args["id"] == "" {
		if id := uuidRe.FindString(message); id != "" {
			args["id"] = id
		}
	}
	if args["status"] == "" {
		if st := findStatus(lower); st != "" {
			args["status"] = st
		}
	}
	return Call{Name: toolPatch, Args: args}
}

func parseFilters(message string) map[string]string {
	args := parseKeyValues(message)
	out := map[string]string{}
	if v := args["status"]; v != "" {
		out["status"] = v
	}
	if v := args["company"]; v != "" {
		out["company"] = v
	}
	if out["status"] == "" {
		if st := findStatus(strings.ToLower(message)); st != "" && strings.Contains(strings.ToLower(message), "status") {
			out["status"] = st
		}
	}
	return out
}

func parseKeyValues(message string) map[string]string {
	out := map[string]string{}
	// company=Acme title="Backend Go" or company=Acme
	re := regexp.MustCompile(`(?i)\b(id|company|title|url|channel|status|notes|contact|location|remote|stack)=("([^"]*)"|'([^']*)'|(\S+))`)
	for _, m := range re.FindAllStringSubmatch(message, -1) {
		key := strings.ToLower(m[1])
		val := m[3] + m[4] + m[5]
		out[key] = strings.TrimSpace(val)
	}
	return out
}

func stripCreateVerb(message string) string {
	re := regexp.MustCompile(`(?i)^\s*(crea(r)?|ajoute(r)?|add|nueva|nouveau|new(\s+application)?)\s*[:\-]?\s*`)
	return strings.TrimSpace(re.ReplaceAllString(message, ""))
}

func splitCompanyTitle(rest string) []string {
	rest = strings.TrimSpace(rest)
	for _, sep := range []string{"|", "/", "—", "–", " - "} {
		if strings.Contains(rest, sep) {
			parts := strings.SplitN(rest, sep, 2)
			c := strings.TrimSpace(parts[0])
			t := strings.TrimSpace(parts[1])
			// drop trailing key=value from title side
			if i := strings.Index(t, " channel="); i > 0 {
				t = strings.TrimSpace(t[:i])
			}
			if c != "" && t != "" && !strings.Contains(c, "=") {
				return []string{c, t}
			}
		}
	}
	return nil
}

func findStatus(lower string) string {
	for st := range models.AllowedStatuses {
		if strings.Contains(lower, st) {
			return st
		}
	}
	return ""
}

func containsAny(s string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}
