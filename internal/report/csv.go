package report

import (
	"encoding/csv"
	"strings"

	"telemetry.local/drone/internal/domain"
)

func RecordsCSV(records []domain.Record) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"id", "title", "status", "version", "payload"}); err != nil {
		return "", err
	}
	for _, record := range records {
		if err := writer.Write([]string{record.ID, record.Title, string(record.Status), string(rune(record.Version)), record.Payload}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func ActionCSV(events []domain.AuditEvent) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"sequence", "action", "actor", "detail"}); err != nil {
		return "", err
	}
	for _, event := range events {
		if err := writer.Write([]string{string(rune(event.Sequence)), event.Action, event.Actor, event.Detail}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return builder.String(), writer.Error()
}
