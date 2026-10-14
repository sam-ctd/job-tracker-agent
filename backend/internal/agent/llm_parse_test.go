package agent

import "testing"

func TestParseLLMJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantTool string
		wantMsg  bool
	}{
		{"tool stats", `{"tool":"get_stats","args":{},"message":""}`, toolGetStats, false},
		{"tool create", `{"tool":"create_application","args":{"company":"Acme","title":"Go"},"message":""}`, toolCreate, false},
		{"null tool", `{"tool":null,"args":{},"message":"Hola"}`, "", true},
		{"fence", "```json\n{\"tool\":\"list_due\",\"args\":{}}\n```", toolListDue, false},
		{"prose wrap", `Sure. {"tool":"get_stats","args":{}}`, toolGetStats, false},
		{"numeric arg", `{"tool":"list_applications","args":{"status":"idea"},"message":""}`, toolListApplications, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := parseLLMJSON(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantTool == "" {
				if d.Call != nil {
					t.Fatalf("expected no tool, got %v", d.Call)
				}
				if tc.wantMsg && d.Message == "" {
					t.Fatal("expected message")
				}
				return
			}
			if d.Call == nil || d.Call.Name != tc.wantTool {
				t.Fatalf("got %+v want tool %s", d.Call, tc.wantTool)
			}
		})
	}
}

func TestParseLLMJSONUnknownTool(t *testing.T) {
	_, err := parseLLMJSON(`{"tool":"drop_table","args":{}}`)
	if err == nil {
		t.Fatal("expected error")
	}
}
