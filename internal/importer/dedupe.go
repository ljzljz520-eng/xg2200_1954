package importer

import (
	"sort"

	"telemetry.local/drone/internal/domain"
)

type DedupeResult struct {
	Unique     []domain.Record
	Duplicates []domain.Record
}

func Deduplicate(records []domain.Record) DedupeResult {
	seen := make(map[string]domain.Record)
	duplicates := make([]domain.Record, 0)
	for _, record := range records {
		if _, ok := seen[record.ID]; ok {
			duplicates = append(duplicates, record)
			continue
		}
		seen[record.ID] = record
	}
	unique := make([]domain.Record, 0, len(seen))
	for _, record := range seen {
		unique = append(unique, record)
	}
	sort.Slice(unique, func(i, j int) bool { return unique[i].ID < unique[j].ID })
	return DedupeResult{Unique: unique, Duplicates: duplicates}
}

func GroupByKind(attachments []domain.Attachment) map[string][]domain.Attachment {
	groups := make(map[string][]domain.Attachment)
	for _, attachment := range attachments {
		groups[attachment.Kind] = append(groups[attachment.Kind], attachment)
	}
	for kind := range groups {
		sort.Slice(groups[kind], func(i, j int) bool { return groups[kind][i].ID < groups[kind][j].ID })
	}
	return groups
}
