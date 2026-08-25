package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Decoder struct {
	Strict bool
}

func (d Decoder) DecodeCSV(r io.Reader, batch string) (Result, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	rows := make([]Row, 0)
	for index := 1; ; index++ {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Result{}, err
		}
		if len(fields) < 4 {
			if d.Strict {
				return Result{}, fmt.Errorf("csv row %d has %d fields", index, len(fields))
			}
			continue
		}
		rows = append(rows, Row{ID: strings.TrimSpace(fields[0]), Title: strings.TrimSpace(fields[1]), Payload: strings.TrimSpace(fields[2]), Kind: strings.TrimSpace(fields[3])})
	}
	return ParseRows(rows, batch)
}

func ValidateRows(rows []Row) []string {
	issues := make([]string, 0)
	for index, row := range rows {
		record := domain.Record{ID: row.ID, Title: row.Title, Payload: row.Payload, Status: domain.StatusPending, Version: 1}
		if err := domain.ValidateRecord(record); err != nil {
			issues = append(issues, fmt.Sprintf("%d:%v", index+1, err))
			continue
		}
		if !domain.AllowedTitle(row.Title) {
			issues = append(issues, fmt.Sprintf("%d:invalid-title", index+1))
		}
	}
	return issues
}
