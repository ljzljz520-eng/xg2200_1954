package flow003

import (
	"strings"
	"testing"

	"telemetry.local/drone/internal/service"
	"telemetry.local/drone/internal/store"
)

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, "importer")
	result, summary, err := app.ImportBatch(strings.NewReader("a|Current Flight Plan|cipher-a|telemetry\nb|expired flight plan|cipher-b|telemetry"), "flow003")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 || len(result.Rejected) != 1 {
		t.Fatalf("result=%+v", result)
	}
	if summary.Accepted != 1 || summary.Rejected != 1 {
		t.Fatalf("summary=%+v", summary)
	}
}
