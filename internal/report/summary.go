package report

import (
	"fmt"
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/importer"
)

type Summary struct {
	Total      int            `json:"total"`
	ByStatus   map[string]int `json:"by_status"`
	Titles     []string       `json:"titles"`
	AuditCount int            `json:"audit_count"`
	Accepted   int            `json:"accepted"`
	Rejected   int            `json:"rejected"`
}

func BuildSummary(records []domain.Record, audits []domain.AuditEvent, imported *importer.Result) Summary {
	summary := Summary{Total: len(records), ByStatus: make(map[string]int), Titles: make([]string, 0, len(records)), AuditCount: len(audits)}
	for _, record := range records {
		summary.ByStatus[string(record.Status)]++
		summary.Titles = append(summary.Titles, record.Title)
	}
	sort.Strings(summary.Titles)
	if imported != nil {
		summary.Accepted = imported.AcceptedCount()
		summary.Rejected = imported.RejectedCount()
	}
	return summary
}

func Render(summary Summary) string {
	statuses := make([]string, 0, len(summary.ByStatus))
	for status := range summary.ByStatus {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	parts := []string{fmt.Sprintf("total=%d", summary.Total), fmt.Sprintf("audits=%d", summary.AuditCount), fmt.Sprintf("accepted=%d", summary.Accepted), fmt.Sprintf("rejected=%d", summary.Rejected)}
	for _, status := range statuses {
		parts = append(parts, fmt.Sprintf("%s=%d", status, summary.ByStatus[status]))
	}
	if len(summary.Titles) > 0 {
		parts = append(parts, "titles="+strings.Join(summary.Titles, ","))
	}
	return strings.Join(parts, " ")
}
