package service

import (
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/report"
)

func FilterAndSort(records []domain.Record, query string) []domain.Record {
	filtered := make([]domain.Record, 0)
	q := strings.ToLower(strings.TrimSpace(query))
	for _, record := range records {
		if q == "" || strings.Contains(strings.ToLower(record.Title), q) || strings.Contains(strings.ToLower(record.ID), q) {
			filtered = append(filtered, record)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Title == filtered[j].Title {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].Title < filtered[j].Title
	})
	return filtered
}

func BuildExport(records []domain.Record, audits []domain.AuditEvent) report.Export {
	return report.NewExport(records, audits)
}

func StatusCounts(records []domain.Record) map[domain.Status]int {
	counts := make(map[domain.Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}
