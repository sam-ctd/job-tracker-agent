package store

import (
	"errors"
	"testing"

	"github.com/sam-ctd/job-tracker-agent/internal/models"
)

func TestValidateFieldsOK(t *testing.T) {
	err := validateFields("Acme", "Stage Go", "https://example.com", "linkedin", "idea", "n", "c", "Lyon", "remote", "Go")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateFieldsRejects(t *testing.T) {
	cases := []struct {
		name string
		fn   func() error
	}{
		{"empty company", func() error {
			return validateFields("", "t", "", "linkedin", "idea", "", "", "", "", "")
		}},
		{"bad url", func() error {
			return validateFields("Acme", "t", "ftp://x", "linkedin", "idea", "", "", "", "", "")
		}},
		{"bad status", func() error {
			return validateFields("Acme", "t", "", "linkedin", "nope", "", "", "", "", "")
		}},
		{"bad channel", func() error {
			return validateFields("Acme", "t", "", "twitter", "idea", "", "", "", "", "")
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestNormalizeCreateDefaults(t *testing.T) {
	in, err := normalizeCreate(models.ApplicationCreate{
		Company: " Acme ",
		Title:   " Stage ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if in.Status != models.StatusIdea || in.Channel != models.ChannelLinkedIn {
		t.Fatalf("unexpected defaults: %+v", in)
	}
	if in.Company != "Acme" || in.Title != "Stage" {
		t.Fatalf("trim failed: %+v", in)
	}
}

func TestParseUUID(t *testing.T) {
	if err := parseUUID("not-a-uuid"); err == nil {
		t.Fatal("expected error")
	}
	if err := parseUUID("9483e199-9460-44ba-b009-9e28c2f5b053"); err != nil {
		t.Fatal(err)
	}
}
