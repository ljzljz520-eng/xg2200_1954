package report

import (
	"encoding/json"
	"sort"

	"telemetry.local/drone/internal/domain"
)

type Export struct {
	Records []domain.Record `json:"records"`
	Audits  []AuditView     `json:"audits"`
}

func NewExport(records []domain.Record, audits []domain.AuditEvent) Export {
	copyRecords := append([]domain.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	return Export{Records: copyRecords, Audits: BuildAuditView(audits)}
}

func (e Export) JSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

func (e Export) Count() int {
	return len(e.Records)
}
