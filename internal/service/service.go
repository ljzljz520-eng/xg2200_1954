package service

import (
	"fmt"
	"io"
	"strings"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/importer"
	"telemetry.local/drone/internal/report"
	"telemetry.local/drone/internal/review"
	"telemetry.local/drone/internal/store"
)

type Service struct {
	Store    *store.Store
	Reviewer review.Reviewer
	Now      int64
}

func New(s *store.Store, actor string) *Service {
	return &Service{Store: s, Reviewer: review.New(actor), Now: 1}
}

func (s *Service) ImportBatch(r io.Reader, batch string) (importer.Result, report.Summary, error) {
	processor := importer.NewProcessor()
	parsed, err := processor.ParseBatch(r, batch, 1)
	if err != nil {
		return importer.Result{}, report.Summary{}, err
	}
	valid := make([]domain.Record, 0, len(parsed.Records))
	for _, record := range parsed.Records {
		if !domain.AllowedTitle(record.Title) {
			parsed.Rejected = append(parsed.Rejected, fmt.Sprintf("%s: invalid title", record.ID))
			continue
		}
		valid = append(valid, record)
	}
	parsed.Records = valid
	for _, record := range parsed.Records {
		record.CreatedAt = s.Now
		record.UpdatedAt = s.Now
		s.Now++
		if err := s.Store.SaveRecord(record); err != nil {
			return importer.Result{}, report.Summary{}, err
		}
		for _, attachment := range parsed.Attachments {
			if attachment.RecordID == record.ID {
				if err := s.Store.SaveAttachment(attachment); err != nil {
					return importer.Result{}, report.Summary{}, err
				}
			}
		}
		workflow := domain.WorkflowFor(record.ID, "lifecycle", 1, "imported")
		if err := s.Store.SaveWorkflow(workflow); err != nil {
			return importer.Result{}, report.Summary{}, err
		}
		builder := domain.NewEventBuilder(record.ID, s.Reviewer.Actor, 2)
		event := builder.Build("import", "record persisted")
		if err := s.Store.SaveAudit(event); err != nil {
			return importer.Result{}, report.Summary{}, err
		}
	}
	records, err := s.Store.ListRecords()
	if err != nil {
		return importer.Result{}, report.Summary{}, err
	}
	audits, err := s.Store.ListAudits("")
	if err != nil {
		return importer.Result{}, report.Summary{}, err
	}
	return parsed, report.BuildSummary(records, audits, &parsed), nil
}

func (s *Service) Search(query string, status domain.Status) ([]domain.Record, error) {
	return s.Store.Search(domain.Filter{Query: query, Status: status})
}

func (s *Service) ReviewAndArchive(id string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	reviewed, event, err := s.Reviewer.Approve(record)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(reviewed); err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveAudit(event); err != nil {
		return domain.Record{}, err
	}
	workflow := domain.WorkflowFor(id, "lifecycle", 3, "approved")
	if err := s.Store.SaveWorkflow(workflow); err != nil {
		return domain.Record{}, err
	}
	archived, archiveEvent, err := s.Reviewer.Archive(reviewed, 3)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(archived); err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveAudit(archiveEvent); err != nil {
		return domain.Record{}, err
	}
	workflow.CurrentStep = 4
	workflow.State = "archived"
	workflow.CompletedAt = s.Now
	if err := s.Store.SaveWorkflow(workflow); err != nil {
		return domain.Record{}, err
	}
	s.Now++
	return archived, nil
}

func (s *Service) PublishTitle(id, title string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Reviewer.ValidateChange(record, title); err != nil {
		return domain.Record{}, err
	}
	events, err := s.Store.ListAudits(id)
	if err != nil {
		return domain.Record{}, err
	}
	sequence := len(events) + 1
	updated, event, err := s.Reviewer.Publish(record, title, sequence)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveAudit(event); err != nil {
		return domain.Record{}, err
	}
	workflow := domain.WorkflowFor(id, "lifecycle", 4, "published")
	if err := s.Store.SaveWorkflow(workflow); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) BuildReport() (report.Summary, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return report.Summary{}, err
	}
	audits, err := s.Store.ListAudits("")
	if err != nil {
		return report.Summary{}, err
	}
	return report.BuildSummary(records, audits, nil), nil
}

func (s *Service) ImportText(text, batch string) (string, error) {
	_, summary, err := s.ImportBatch(strings.NewReader(text), batch)
	if err != nil {
		return "", err
	}
	return report.Render(summary), nil
}

func (s *Service) RequireRecord(id string) (domain.Record, error) {
	r, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, fmt.Errorf("record %s: %w", id, err)
	}
	return r, nil
}
