package service

import (
	"fmt"
	"sort"

	"telemetry.local/drone/internal/domain"
)

type ReconcileIssue struct {
	RecordID string
	Code     string
	Detail   string
}

func Reconcile(records []domain.Record, audits []domain.AuditEvent) []ReconcileIssue {
	issues := make([]ReconcileIssue, 0)
	byRecord := make(map[string][]domain.AuditEvent)
	for _, event := range audits {
		byRecord[event.RecordID] = append(byRecord[event.RecordID], event)
	}
	for _, record := range records {
		events := byRecord[record.ID]
		if len(events) == 0 {
			issues = append(issues, ReconcileIssue{RecordID: record.ID, Code: "missing-audit", Detail: "record has no audit events"})
			continue
		}
		for index, event := range events {
			if event.Sequence != index+1 {
				issues = append(issues, ReconcileIssue{RecordID: record.ID, Code: "sequence-gap", Detail: fmt.Sprintf("expected %d got %d", index+1, event.Sequence)})
			}
		}
		if record.Status == domain.StatusArchived && events[len(events)-1].Action != "archive" {
			issues = append(issues, ReconcileIssue{RecordID: record.ID, Code: "archive-mismatch", Detail: "archived record has a different final action"})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].RecordID < issues[j].RecordID || (issues[i].RecordID == issues[j].RecordID && issues[i].Code < issues[j].Code)
	})
	return issues
}

func (s *Service) ReconcileStore() ([]ReconcileIssue, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	audits, err := s.Store.ListAudits("")
	if err != nil {
		return nil, err
	}
	return Reconcile(records, audits), nil
}
