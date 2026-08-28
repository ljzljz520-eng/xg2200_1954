package service

import (
	"fmt"
	"sort"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/review"
)

type LifecycleResult struct {
	Record  domain.Record
	Events  []domain.AuditEvent
	State   string
	Message string
}

func (s *Service) Advance(recordID string, actions []string) (LifecycleResult, error) {
	record, err := s.Store.LoadRecord(recordID)
	if err != nil {
		return LifecycleResult{}, err
	}
	events, err := s.Store.ListAudits(recordID)
	if err != nil {
		return LifecycleResult{}, err
	}
	for _, action := range actions {
		decision := review.Check(record, action)
		if !decision.Allowed {
			return LifecycleResult{}, fmt.Errorf("%s: %s", action, decision.Reason)
		}
		switch action {
		case "approve":
			record, _, err = s.Reviewer.Approve(record)
		case "publish":
			record, _, err = s.Reviewer.Publish(record, record.Title+" updated", len(events)+1)
		case "archive":
			record, _, err = s.Reviewer.Archive(record, len(events)+1)
		}
		if err != nil {
			return LifecycleResult{}, err
		}
		events = append(events, domain.AuditEvent{ID: domain.AuditID(record.ID, len(events)+1), RecordID: record.ID, Action: action, Actor: s.Reviewer.Actor, Detail: "lifecycle advance", Sequence: len(events) + 1})
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return LifecycleResult{}, err
	}
	for _, event := range events {
		if err := s.Store.SaveAudit(event); err != nil {
			return LifecycleResult{}, err
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].Sequence < events[j].Sequence })
	return LifecycleResult{Record: record, Events: events, State: string(record.Status), Message: "lifecycle advanced"}, nil
}

func (s *Service) PendingQueue() ([]string, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	queue := review.NewQueue(records)
	ids := make([]string, 0, queue.Len())
	for {
		record, ok := queue.Next()
		if !ok {
			break
		}
		ids = append(ids, record.ID)
	}
	return ids, nil
}
