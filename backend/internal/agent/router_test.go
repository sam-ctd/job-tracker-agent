package agent

import "testing"

func TestDecideTool(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"muéstrame las stats", toolGetStats},
		{"combien de candidatures ?", toolGetStats},
		{"lista las candidaturas", toolListApplications},
		{"show applications", toolListApplications},
		{"hola qué tal", ""},
		{"", ""},
	}
	for _, tc := range cases {
		got := decideTool(tc.in)
		if got != tc.want {
			t.Fatalf("decideTool(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}
