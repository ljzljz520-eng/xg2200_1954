package store

import (
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Query struct {
	Text         string
	Statuses     []domain.Status
	IncludeAudit bool
	Limit        int
}

type QueryResult struct {
	Records []domain.Record
	Audits  []domain.AuditEvent
}

func (s *Store) ExecuteQuery(query Query) (QueryResult, error) {
	records, err := s.ListRecords()
	if err != nil {
		return QueryResult{}, err
	}
	allowed := make(map[domain.Status]bool)
	for _, status := range query.Statuses {
		allowed[status] = true
	}
	filtered := make([]domain.Record, 0, len(records))
	needle := strings.ToLower(strings.TrimSpace(query.Text))
	for _, record := range records {
		if len(allowed) > 0 && !allowed[record.Status] {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(record.Title), needle) && !strings.Contains(strings.ToLower(record.Payload), needle) {
			continue
		}
		filtered = append(filtered, record)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].UpdatedAt > filtered[j].UpdatedAt || (filtered[i].UpdatedAt == filtered[j].UpdatedAt && filtered[i].ID < filtered[j].ID)
	})
	if query.Limit > 0 && len(filtered) > query.Limit {
		filtered = filtered[:query.Limit]
	}
	result := QueryResult{Records: filtered}
	if query.IncludeAudit {
		result.Audits, err = s.ListAudits("")
	}
	return result, err
}

func (s *Store) CountByStatus() (map[domain.Status]int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	counts := make(map[domain.Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts, nil
}
