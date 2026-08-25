package flow012

import (
	"strings"
	"testing"

	"telemetry.local/drone/internal/importer"
)

func Test2200BusinessRegression(t *testing.T) {
	processor := importer.NewProcessor()
	first, err := processor.ParseBatch(strings.NewReader("alpha|Current Flight Plan|cipher-a|telemetry"), "first", 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := processor.ParseBatch(strings.NewReader("delta|Expired Flight Plan|cipher-d|telemetry"), "second", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Records) != 1 || len(second.Records) != 1 {
		t.Fatalf("records first=%d second=%d", len(first.Records), len(second.Records))
	}
	if first.Records[0].Title != "Current Flight Plan" {
		t.Fatalf("first record title=%q", first.Records[0].Title)
	}
}
