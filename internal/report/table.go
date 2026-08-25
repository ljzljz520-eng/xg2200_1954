package report

import (
	"fmt"
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

func RenderRecords(records []domain.Record) string {
	copyRecords := append([]domain.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	lines := []string{"ID | STATUS | VERSION | TITLE"}
	for _, record := range copyRecords {
		lines = append(lines, fmt.Sprintf("%s | %s | %d | %s", record.ID, record.Status, record.Version, record.Title))
	}
	return strings.Join(lines, "\n")
}

func RenderTimeline(events []domain.AuditEvent) string {
	views := BuildAuditView(events)
	lines := make([]string, 0, len(views))
	for _, event := range views {
		lines = append(lines, fmt.Sprintf("%d %s %s %s", event.Sequence, event.Action, event.Actor, event.Detail))
	}
	return strings.Join(lines, "\n")
}
