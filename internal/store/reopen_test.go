package store

import (
	"testing"

	"telemetry.local/drone/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/reopen.db"
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.Record{ID: "persisted", Title: "Current Flight Plan", Status: domain.StatusApproved, Payload: "cipher", Version: 1, CreatedAt: 1, UpdatedAt: 1}
	if err := db.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAudit(domain.AuditEvent{ID: "audit", RecordID: record.ID, Action: "import", Actor: "test", Detail: "saved", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveWorkflow(domain.WorkflowFor(record.ID, "lifecycle", 2, "approved")); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAttachment(domain.Attachment{ID: "attachment", RecordID: record.ID, Name: "telemetry.bin", Digest: "digest", Size: 6, Kind: "telemetry"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	loaded, err := db.LoadRecord(record.ID)
	if err != nil || loaded.Title != record.Title {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	if events, _ := db.ListAudits(record.ID); len(events) != 1 {
		t.Fatalf("events=%d", len(events))
	}
	if attachments, _ := db.ListAttachments(record.ID); len(attachments) != 1 {
		t.Fatalf("attachments=%d", len(attachments))
	}
}
