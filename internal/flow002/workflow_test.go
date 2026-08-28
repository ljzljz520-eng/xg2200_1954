package flow002

import (
	"strings"
	"testing"

	"telemetry.local/drone/internal/service"
	"telemetry.local/drone/internal/store"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, "reviewer-b")
	result, _, err := app.ImportBatch(strings.NewReader("b|Night Survey Route|cipher-b|telemetry"), "flow002")
	if err != nil {
		t.Fatal(err)
	}
	found, err := app.Search("Night", "")
	if err != nil || len(found) != 1 || found[0].ID != result.Records[0].ID {
		t.Fatalf("found=%+v err=%v", found, err)
	}
}

func TestWorkflowSearchUpdatePublishActive(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, "reviewer-b")
	result, _, err := app.ImportBatch(strings.NewReader("b|Night Survey Route|cipher-b|telemetry"), "flow002")
	if err != nil {
		t.Fatal(err)
	}
	approved, _, err := app.Reviewer.Approve(result.Records[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := db.SaveRecord(approved); err != nil {
		t.Fatal(err)
	}
	updated, err := app.PublishTitle(result.Records[0].ID, "Day Survey Route")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Day Survey Route" || updated.Version < 2 {
		t.Fatalf("updated=%+v", updated)
	}
}
