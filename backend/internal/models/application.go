package models

import "time"

const (
	StatusIdea      = "idea"
	StatusEnviada   = "enviada"
	StatusRelance   = "relance"
	StatusEntretien = "entretien"
	StatusOffre     = "offre"
	StatusRefus     = "refus"

	ChannelLinkedIn       = "linkedin"
	ChannelWTTJ           = "wttj"
	ChannelIndeed         = "indeed"
	ChannelHelloWork      = "hellowork"
	ChannelAPEC           = "apec"
	ChannelJobTeaser      = "jobteaser"
	ChannelLesJeudis      = "lesjeudis"
	ChannelChooseYourBoss = "chooseyourboss"
	ChannelAlumniESTIA    = "alumni_estia"
)

var AllowedStatuses = map[string]struct{}{
	StatusIdea: {}, StatusEnviada: {}, StatusRelance: {},
	StatusEntretien: {}, StatusOffre: {}, StatusRefus: {},
}

var AllowedChannels = map[string]struct{}{
	ChannelLinkedIn: {}, ChannelWTTJ: {}, ChannelIndeed: {},
	ChannelHelloWork: {}, ChannelAPEC: {}, ChannelJobTeaser: {},
	ChannelLesJeudis: {}, ChannelChooseYourBoss: {}, ChannelAlumniESTIA: {},
}

type Application struct {
	ID           string     `json:"id"`
	Company      string     `json:"company"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	Channel      string     `json:"channel"`
	Status       string     `json:"status"`
	SentAt       *time.Time `json:"sent_at"`
	NextFollowUp *time.Time `json:"next_follow_up"`
	Notes        string     `json:"notes"`
	Contact      string     `json:"contact"`
	Location     string     `json:"location"`
	Remote       string     `json:"remote"`
	Stack        string     `json:"stack"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ApplicationCreate struct {
	Company      string     `json:"company"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	Channel      string     `json:"channel"`
	Status       string     `json:"status"`
	SentAt       *time.Time `json:"sent_at"`
	NextFollowUp *time.Time `json:"next_follow_up"`
	Notes        string     `json:"notes"`
	Contact      string     `json:"contact"`
	Location     string     `json:"location"`
	Remote       string     `json:"remote"`
	Stack        string     `json:"stack"`
}

type ApplicationPatch struct {
	Company       *string    `json:"company"`
	Title         *string    `json:"title"`
	URL           *string    `json:"url"`
	Channel       *string    `json:"channel"`
	Status        *string    `json:"status"`
	SentAt        *time.Time `json:"sent_at"`
	NextFollowUp  *time.Time `json:"next_follow_up"`
	ClearSentAt   bool       `json:"clear_sent_at"`
	ClearFollowUp bool       `json:"clear_next_follow_up"`
	Notes         *string    `json:"notes"`
	Contact       *string    `json:"contact"`
	Location      *string    `json:"location"`
	Remote        *string    `json:"remote"`
	Stack         *string    `json:"stack"`
}

type Stats struct {
	ByStatus  map[string]int `json:"by_status"`
	ByChannel map[string]int `json:"by_channel"`
	Total     int            `json:"total"`
}
