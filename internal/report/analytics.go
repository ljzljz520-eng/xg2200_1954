package report

import (
	"sort"

	"telemetry.local/drone/internal/domain"
)

type StatusRow struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Share  int    `json:"share_percent"`
}

func StatusTable(records []domain.Record) []StatusRow {
	counts := make(map[string]int)
	for _, record := range records {
		counts[string(record.Status)]++
	}
	statuses := make([]string, 0, len(counts))
	for status := range counts {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	rows := make([]StatusRow, 0, len(statuses))
	for _, status := range statuses {
		share := 0
		if len(records) > 0 {
			share = counts[status] * 100 / len(records)
		}
		rows = append(rows, StatusRow{Status: status, Count: counts[status], Share: share})
	}
	return rows
}

func StaleRecords(records []domain.Record, cutoff int64) []domain.Record {
	stale := make([]domain.Record, 0)
	for _, record := range records {
		if record.UpdatedAt < cutoff && record.Status != domain.StatusArchived {
			stale = append(stale, record)
		}
	}
	sort.Slice(stale, func(i, j int) bool { return stale[i].UpdatedAt < stale[j].UpdatedAt })
	return stale
}

func VersionHistogram(records []domain.Record) map[int]int {
	histogram := make(map[int]int)
	for _, record := range records {
		histogram[record.Version]++
	}
	return histogram
}
