package importer

import (
	"strings"
	"testing"
)

func TestParseAndValidateBatch(t *testing.T) {
	result, err := ParseBatch(strings.NewReader(DeterministicBatch()), "fixture", 1)
	if err != nil {
		t.Fatal(err)
	}
	summary := ValidateBatch(result)
	if summary.Accepted != 2 || summary.Rejected != 1 {
		t.Fatalf("summary=%+v", summary)
	}
}
