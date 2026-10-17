package store

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sam-ctd/job-tracker-agent/internal/models"
)

const appColumns = `id, company, title, url, channel, status, priority, sent_at, next_follow_up, notes, contact, location, remote, stack, created_at, updated_at`

const priorityOrder = `CASE priority WHEN 'alta' THEN 0 WHEN 'media' THEN 1 WHEN 'baja' THEN 2 ELSE 3 END`

func (s *Store) List(ctx context.Context, status, company string) ([]models.Application, error) {
	q := `SELECT ` + appColumns + ` FROM applications WHERE 1=1`
	args := []any{}
	n := 1
	if status != "" {
		q += fmt.Sprintf(` AND status=$%d`, n)
		args = append(args, status)
		n++
	}
	if company != "" {
		q += fmt.Sprintf(` AND company ILIKE $%d`, n)
		args = append(args, "%"+company+"%")
		n++
	}
	q += ` ORDER BY ` + priorityOrder + `, updated_at DESC`

	rows, err := s.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApps(rows)
}

func (s *Store) Get(ctx context.Context, id string) (models.Application, error) {
	if err := parseUUID(id); err != nil {
		return models.Application{}, fmt.Errorf("%w: id", ErrValidation)
	}
	row := s.pool.QueryRow(ctx,
		`SELECT `+appColumns+` FROM applications WHERE id=$1`, id)
	return scanApp(row)
}

func (s *Store) Create(ctx context.Context, in models.ApplicationCreate) (models.Application, error) {
	in, err := normalizeCreate(in)
	if err != nil {
		return models.Application{}, err
	}
	row := s.pool.QueryRow(ctx, `
		INSERT INTO applications (company, title, url, channel, status, priority, sent_at, next_follow_up, notes, contact, location, remote, stack)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING `+appColumns,
		in.Company, in.Title, in.URL, in.Channel, in.Status, in.Priority, in.SentAt, in.NextFollowUp,
		in.Notes, in.Contact, in.Location, in.Remote, in.Stack,
	)
	return scanApp(row)
}

func (s *Store) Patch(ctx context.Context, id string, in models.ApplicationPatch) (models.Application, error) {
	if err := parseUUID(id); err != nil {
		return models.Application{}, fmt.Errorf("%w: id", ErrValidation)
	}
	cur, err := s.Get(ctx, id)
	if err != nil {
		return models.Application{}, err
	}
	merged, err := applyPatch(cur, in)
	if err != nil {
		return models.Application{}, err
	}
	row := s.pool.QueryRow(ctx, `
		UPDATE applications SET
			company=$2, title=$3, url=$4, channel=$5, status=$6, priority=$7,
			sent_at=$8, next_follow_up=$9, notes=$10, contact=$11,
			location=$12, remote=$13, stack=$14, updated_at=NOW()
		WHERE id=$1
		RETURNING `+appColumns,
		id, merged.Company, merged.Title, merged.URL, merged.Channel, merged.Status, merged.Priority,
		merged.SentAt, merged.NextFollowUp, merged.Notes, merged.Contact,
		merged.Location, merged.Remote, merged.Stack,
	)
	return scanApp(row)
}

func (s *Store) Delete(ctx context.Context, id string) error {
	if err := parseUUID(id); err != nil {
		return fmt.Errorf("%w: id", ErrValidation)
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM applications WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) Due(ctx context.Context) ([]models.Application, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+appColumns+`
		FROM applications
		WHERE next_follow_up IS NOT NULL AND next_follow_up <= NOW()
		  AND status NOT IN ('offre', 'refus')
		ORDER BY `+priorityOrder+`, next_follow_up ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApps(rows)
}

func (s *Store) Stats(ctx context.Context) (models.Stats, error) {
	stats := models.Stats{
		ByStatus:  map[string]int{},
		ByChannel: map[string]int{},
	}
	rows, err := s.pool.Query(ctx, `SELECT status, COUNT(*) FROM applications GROUP BY status`)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByStatus[status] = n
		stats.Total += n
	}
	rows.Close()

	rows, err = s.pool.Query(ctx, `SELECT channel, COUNT(*) FROM applications GROUP BY channel`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var channel string
		var n int
		if err := rows.Scan(&channel, &n); err != nil {
			return stats, err
		}
		stats.ByChannel[channel] = n
	}
	return stats, rows.Err()
}

type scannable interface {
	Scan(dest ...any) error
}

func scanApp(row scannable) (models.Application, error) {
	var a models.Application
	err := row.Scan(
		&a.ID, &a.Company, &a.Title, &a.URL, &a.Channel, &a.Status, &a.Priority,
		&a.SentAt, &a.NextFollowUp, &a.Notes, &a.Contact,
		&a.Location, &a.Remote, &a.Stack, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return a, ErrNotFound
		}
		return a, err
	}
	return a, nil
}

func scanApps(rows pgx.Rows) ([]models.Application, error) {
	out := []models.Application{}
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func parseUUID(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func normalizeCreate(in models.ApplicationCreate) (models.ApplicationCreate, error) {
	in.Company = strings.TrimSpace(in.Company)
	in.Title = strings.TrimSpace(in.Title)
	in.URL = strings.TrimSpace(in.URL)
	in.Channel = strings.TrimSpace(in.Channel)
	in.Status = strings.TrimSpace(in.Status)
	in.Priority = strings.TrimSpace(in.Priority)
	in.Notes = strings.TrimSpace(in.Notes)
	in.Contact = strings.TrimSpace(in.Contact)
	in.Location = strings.TrimSpace(in.Location)
	in.Remote = strings.TrimSpace(in.Remote)
	in.Stack = strings.TrimSpace(in.Stack)

	if in.Status == "" {
		in.Status = models.StatusIdea
	}
	if in.Channel == "" {
		in.Channel = models.ChannelLinkedIn
	}
	if in.Priority == "" {
		in.Priority = models.PriorityMedia
	}
	if err := validateFields(in.Company, in.Title, in.URL, in.Channel, in.Status, in.Priority, in.Notes, in.Contact, in.Location, in.Remote, in.Stack); err != nil {
		return in, err
	}
	return in, nil
}

func applyPatch(cur models.Application, in models.ApplicationPatch) (models.Application, error) {
	if in.Company != nil {
		cur.Company = strings.TrimSpace(*in.Company)
	}
	if in.Title != nil {
		cur.Title = strings.TrimSpace(*in.Title)
	}
	if in.URL != nil {
		cur.URL = strings.TrimSpace(*in.URL)
	}
	if in.Channel != nil {
		cur.Channel = strings.TrimSpace(*in.Channel)
	}
	if in.Status != nil {
		cur.Status = strings.TrimSpace(*in.Status)
	}
	if in.Priority != nil {
		cur.Priority = strings.TrimSpace(*in.Priority)
	}
	if in.ClearSentAt {
		cur.SentAt = nil
	} else if in.SentAt != nil {
		cur.SentAt = in.SentAt
	}
	if in.ClearFollowUp {
		cur.NextFollowUp = nil
	} else if in.NextFollowUp != nil {
		cur.NextFollowUp = in.NextFollowUp
	}
	if in.Notes != nil {
		cur.Notes = strings.TrimSpace(*in.Notes)
	}
	if in.Contact != nil {
		cur.Contact = strings.TrimSpace(*in.Contact)
	}
	if in.Location != nil {
		cur.Location = strings.TrimSpace(*in.Location)
	}
	if in.Remote != nil {
		cur.Remote = strings.TrimSpace(*in.Remote)
	}
	if in.Stack != nil {
		cur.Stack = strings.TrimSpace(*in.Stack)
	}
	if cur.Priority == "" {
		cur.Priority = models.PriorityMedia
	}
	if err := validateFields(cur.Company, cur.Title, cur.URL, cur.Channel, cur.Status, cur.Priority, cur.Notes, cur.Contact, cur.Location, cur.Remote, cur.Stack); err != nil {
		return cur, err
	}
	return cur, nil
}

func validateFields(company, title, urlStr, channel, status, priority, notes, contact, location, remote, stack string) error {
	if company == "" || len(company) > 200 {
		return fmt.Errorf("%w: company", ErrValidation)
	}
	if title == "" || len(title) > 300 {
		return fmt.Errorf("%w: title", ErrValidation)
	}
	if len(urlStr) > 2000 {
		return fmt.Errorf("%w: url", ErrValidation)
	}
	if urlStr != "" {
		u, err := url.ParseRequestURI(urlStr)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("%w: url", ErrValidation)
		}
	}
	if _, ok := models.AllowedChannels[channel]; !ok {
		return fmt.Errorf("%w: channel", ErrValidation)
	}
	if _, ok := models.AllowedStatuses[status]; !ok {
		return fmt.Errorf("%w: status", ErrValidation)
	}
	if _, ok := models.AllowedPriorities[priority]; !ok {
		return fmt.Errorf("%w: priority", ErrValidation)
	}
	if len(notes) > 5000 {
		return fmt.Errorf("%w: notes", ErrValidation)
	}
	if len(contact) > 500 {
		return fmt.Errorf("%w: contact", ErrValidation)
	}
	if len(location) > 300 {
		return fmt.Errorf("%w: location", ErrValidation)
	}
	if len(remote) > 100 {
		return fmt.Errorf("%w: remote", ErrValidation)
	}
	if len(stack) > 500 {
		return fmt.Errorf("%w: stack", ErrValidation)
	}
	return nil
}
