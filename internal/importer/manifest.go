package importer

import (
	"fmt"
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Manifest struct {
	BatchID string
	Rows    int
	Kinds   []string
	Digest  string
}

func BuildManifest(batch string, result Result) Manifest {
	kinds := make(map[string]bool)
	for _, attachment := range result.Attachments {
		kinds[attachment.Kind] = true
	}
	list := make([]string, 0, len(kinds))
	for kind := range kinds {
		list = append(list, kind)
	}
	sort.Strings(list)
	parts := []string{batch, fmt.Sprint(len(result.Records)), strings.Join(list, ",")}
	for _, record := range result.Records {
		parts = append(parts, record.ID, record.Title)
	}
	return Manifest{BatchID: batch, Rows: len(result.Records) + len(result.Rejected), Kinds: list, Digest: domain.StableID(parts...)}
}

func CompareManifest(left, right Manifest) bool {
	if left.BatchID != right.BatchID || left.Rows != right.Rows || left.Digest != right.Digest {
		return false
	}
	if len(left.Kinds) != len(right.Kinds) {
		return false
	}
	for index := range left.Kinds {
		if left.Kinds[index] != right.Kinds[index] {
			return false
		}
	}
	return true
}
