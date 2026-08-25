package review

import (
	"fmt"

	"telemetry.local/drone/internal/domain"
)

type Reviewer struct {
	Actor string
}

func New(actor string) Reviewer {
	if actor == "" {
		actor = "system-reviewer"
	}
	return Reviewer{Actor: actor}
}

func (r Reviewer) Approve(record domain.Record) (domain.Record, domain.AuditEvent, error) {
	if err := domain.Transition(record.Status, domain.StatusApproved); err != nil {
		return record, domain.AuditEvent{}, err
	}
	if !domain.AllowedTitle(record.Title) {
		return record, domain.AuditEvent{}, fmt.Errorf("title is not eligible for approval")
	}
	record.Status = domain.StatusApproved
	record.Version++
	builder := domain.NewEventBuilder(record.ID, r.Actor, 1)
	event := builder.Build("approve", "review checks passed")
	return record, event, nil
}

func (r Reviewer) Publish(record domain.Record, title string, sequence int) (domain.Record, domain.AuditEvent, error) {
	if record.Status != domain.StatusApproved {
		return record, domain.AuditEvent{}, fmt.Errorf("record must be approved before publishing")
	}
	newTitle := domain.NormalizeTitle(title)
	if !domain.AllowedTitle(newTitle) {
		return record, domain.AuditEvent{}, fmt.Errorf("title is not eligible for publishing")
	}
	if newTitle == record.Title {
		builder := domain.NewEventBuilder(record.ID, r.Actor, sequence)
		return record, builder.Build("publish-noop", "title unchanged"), nil
	}
	record.Title = newTitle
	record.Status = domain.StatusPublished
	record.Version++
	builder := domain.NewEventBuilder(record.ID, r.Actor, sequence)
	event := builder.Build("publish", fmt.Sprintf("version=%d", record.Version))
	return record, event, nil
}

func (r Reviewer) Archive(record domain.Record, sequence int) (domain.Record, domain.AuditEvent, error) {
	if err := domain.Transition(record.Status, domain.StatusArchived); err != nil {
		return record, domain.AuditEvent{}, err
	}
	record.Status = domain.StatusArchived
	record.Version++
	builder := domain.NewEventBuilder(record.ID, r.Actor, sequence)
	event := builder.Build("archive", "retained audit and attachments")
	return record, event, nil
}

func (r Reviewer) ValidateChange(record domain.Record, title string) error {
	if record.Status == domain.StatusArchived {
		return fmt.Errorf("archived record cannot change")
	}
	if !domain.AllowedTitle(title) {
		return fmt.Errorf("invalid title")
	}
	return nil
}
