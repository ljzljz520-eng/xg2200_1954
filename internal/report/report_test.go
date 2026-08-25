package report

import (
	"testing"

	"telemetry.local/drone/internal/domain"
)

func TestSummaryRender(t *testing.T) {
	records := []domain.Record{{ID: "b", Title: "B", Status: domain.StatusApproved}, {ID: "a", Title: "A", Status: domain.StatusPending}}
	summary := BuildSummary(records, []domain.AuditEvent{{Action: "import"}}, nil)
	text := Render(summary)
	if summary.Total != 2 || summary.ByStatus["approved"] != 1 || text == "" {
		t.Fatalf("summary=%+v text=%q", summary, text)
	}
}
