package importer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Row struct {
	ID      string
	Title   string
	Payload string
	Kind    string
}

type Result struct {
	Records     []domain.Record
	Attachments []domain.Attachment
	Rejected    []string
}

func (r Result) AcceptedCount() int {
	return len(r.Records)
}

func (r Result) RejectedCount() int {
	return len(r.Rejected)
}

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func parseLine(line string) (Row, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return Row{}, fmt.Errorf("expected four fields, got %d", len(parts))
	}
	if strings.TrimSpace(parts[0]) == "" {
		return Row{}, fmt.Errorf("missing id")
	}
	return Row{ID: strings.TrimSpace(parts[0]), Title: strings.TrimSpace(parts[1]), Payload: strings.TrimSpace(parts[2]), Kind: strings.TrimSpace(parts[3])}, nil
}

func ParseBatch(r io.Reader, batch string, baseRow int) (Result, error) {
	return NewProcessor().ParseBatch(r, batch, baseRow)
}

func (p *Processor) ParseBatch(r io.Reader, batch string, baseRow int) (Result, error) {
	result := Result{Records: make([]domain.Record, 0), Attachments: make([]domain.Attachment, 0), Rejected: make([]string, 0)}
	scanner := bufio.NewScanner(r)
	row := baseRow
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parsed, err := parseLine(line)
		if err != nil {
			result.Rejected = append(result.Rejected, fmt.Sprintf("row %d: %v", row, err))
			row++
			continue
		}
		record := domain.Record{ID: parsed.ID, Title: parsed.Title, Status: domain.StatusPending, Payload: parsed.Payload, Version: 1, CreatedAt: int64(row), UpdatedAt: int64(row)}
		if err := domain.ValidateRecord(record); err != nil {
			result.Rejected = append(result.Rejected, fmt.Sprintf("row %d: %v", row, err))
			row++
			continue
		}
		attachment := domain.Attachment{ID: domain.AttachmentID(record.ID, parsed.Kind), RecordID: record.ID, Name: parsed.Kind + ".bin", Digest: domain.StableID(record.Payload), Size: len(record.Payload), Kind: parsed.Kind}
		result.Records = append(result.Records, record)
		result.Attachments = append(result.Attachments, attachment)
		row++
	}
	if err := scanner.Err(); err != nil {
		return Result{}, err
	}
	return result, nil
}

func ParseRows(rows []Row, batch string) (Result, error) {
	var b strings.Builder
	for i, row := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(row.ID)
		b.WriteByte('|')
		b.WriteString(row.Title)
		b.WriteByte('|')
		b.WriteString(row.Payload)
		b.WriteByte('|')
		b.WriteString(row.Kind)
	}
	return ParseBatch(strings.NewReader(b.String()), batch, 1)
}

func ParseCount(value string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(value))
}
