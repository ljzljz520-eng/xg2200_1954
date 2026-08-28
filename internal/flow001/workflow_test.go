package flow001

import (
	"strings"
	"testing"

	"telemetry.local/drone/internal/service"
	"telemetry.local/drone/internal/store"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	path := t.TempDir() + "/flow.db"
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, "reviewer-a")
	result, _, err := app.ImportBatch(strings.NewReader("a|Current Flight Plan|cipher-a|telemetry"), "flow001")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("records=%d", len(result.Records))
	}
	archived, err := app.ReviewAndArchive(result.Records[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != "archived" {
		t.Fatalf("status=%s", archived.Status)
	}
	events, err := db.ListAudits(archived.ID)
	if err != nil || len(events) != 3 {
		t.Fatalf("events=%d err=%v", len(events), err)
	}
}
