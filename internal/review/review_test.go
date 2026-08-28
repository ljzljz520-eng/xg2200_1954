package review

import (
	"testing"

	"telemetry.local/drone/internal/domain"
)

func TestReviewerTransitions(t *testing.T) {
	reviewer := New("qa")
	record := domain.Record{ID: "r", Title: "Current Flight Plan", Status: domain.StatusPending, Payload: "cipher", Version: 1}
	approved, event, err := reviewer.Approve(record)
	if err != nil || approved.Status != domain.StatusApproved || event.Action != "approve" {
		t.Fatalf("approved=%+v event=%+v err=%v", approved, event, err)
	}
	published, _, err := reviewer.Publish(approved, "Day Flight Plan", 2)
	if err != nil || published.Status != domain.StatusPublished {
		t.Fatalf("published=%+v err=%v", published, err)
	}
}
