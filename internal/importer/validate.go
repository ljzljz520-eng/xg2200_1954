package importer

import (
	"fmt"
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type ValidationSummary struct {
	Accepted int
	Rejected int
	Reasons  []string
}

func ValidateBatch(result Result) ValidationSummary {
	summary := ValidationSummary{Accepted: len(result.Records), Rejected: len(result.Rejected), Reasons: append([]string(nil), result.Rejected...)}
	for _, record := range result.Records {
		if err := domain.ValidateRecord(record); err != nil {
			summary.Rejected++
			summary.Accepted--
			summary.Reasons = append(summary.Reasons, fmt.Sprintf("%s: %v", record.ID, err))
		} else if !domain.AllowedTitle(record.Title) {
			summary.Rejected++
			summary.Accepted--
			summary.Reasons = append(summary.Reasons, fmt.Sprintf("%s: invalid title", record.ID))
		}
	}
	sort.Strings(summary.Reasons)
	return summary
}

func IsDuplicate(records []domain.Record) []string {
	seen := make(map[string]bool)
	duplicates := make([]string, 0)
	for _, record := range records {
		if seen[record.ID] {
			duplicates = append(duplicates, record.ID)
		}
		seen[record.ID] = true
	}
	sort.Strings(duplicates)
	return duplicates
}

func Explain(result Result) string {
	summary := ValidateBatch(result)
	if summary.Rejected == 0 {
		return fmt.Sprintf("accepted=%d rejected=0", summary.Accepted)
	}
	return fmt.Sprintf("accepted=%d rejected=%d reasons=%s", summary.Accepted, summary.Rejected, strings.Join(summary.Reasons, ";"))
}
