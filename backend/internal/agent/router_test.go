package agent

import "testing"

func TestDecide(t *testing.T) {
	cases := []struct {
		in       string
		wantTool string
		wantCo   string
		wantTi   string
		wantSt   string
		wantID   bool
	}{
		{"muéstrame las stats", toolGetStats, "", "", "", false},
		{"lista las candidaturas", toolListApplications, "", "", "", false},
		{"qué hay due", toolListDue, "", "", "", false},
		{"muestra relances pendientes", toolListDue, "", "", "", false},
		{"crea Acme | Backend Go", toolCreate, "Acme", "Backend Go", "", false},
		{`crear company=Estia title="DevOps"`, toolCreate, "Estia", "DevOps", "", false},
		{"patch id=11111111-1111-1111-1111-111111111111 status=enviada", toolPatch, "", "", "enviada", true},
		{"actualiza 22222222-2222-2222-2222-222222222222 a entretien", toolPatch, "", "", "entretien", true},
		{"hola", "", "", "", "", false},
	}
	for _, tc := range cases {
		call, ok := decide(tc.in)
		if tc.wantTool == "" {
			if ok {
				t.Fatalf("decide(%q) unexpected tool %q", tc.in, call.Name)
			}
			continue
		}
		if !ok || call.Name != tc.wantTool {
			t.Fatalf("decide(%q)=(%q,%v) want tool %q", tc.in, call.Name, ok, tc.wantTool)
		}
		if tc.wantCo != "" && call.Args["company"] != tc.wantCo {
			t.Fatalf("decide(%q) company=%q want %q", tc.in, call.Args["company"], tc.wantCo)
		}
		if tc.wantTi != "" && call.Args["title"] != tc.wantTi {
			t.Fatalf("decide(%q) title=%q want %q", tc.in, call.Args["title"], tc.wantTi)
		}
		if tc.wantSt != "" && call.Args["status"] != tc.wantSt {
			t.Fatalf("decide(%q) status=%q want %q", tc.in, call.Args["status"], tc.wantSt)
		}
		if tc.wantID && call.Args["id"] == "" {
			t.Fatalf("decide(%q) missing id", tc.in)
		}
	}
}

func TestSystemPromptHasTools(t *testing.T) {
	p := SystemPrompt()
	for _, name := range []string{toolGetStats, toolCreate, toolPatch, toolListDue} {
		if !containsAny(p, name) {
			t.Fatalf("SystemPrompt missing %s", name)
		}
	}
}

func TestCatalogCount(t *testing.T) {
	if n := len(Catalog()); n != 5 {
		t.Fatalf("Catalog len=%d want 5", n)
	}
}
