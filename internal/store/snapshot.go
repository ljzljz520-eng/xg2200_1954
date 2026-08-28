package store

import (
	"encoding/json"
	"sort"

	"telemetry.local/drone/internal/domain"
)

type Snapshot struct {
	Records     []domain.Record
	Audits      []domain.AuditEvent
	Workflows   []domain.Workflow
	Attachments []domain.Attachment
}

func (s *Store) ExportSnapshot() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAudits("")
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return Snapshot{}, err
	}
	workflows := make([]domain.Workflow, 0)
	for _, record := range records {
		w, loadErr := s.LoadWorkflow(domain.WorkflowID(record.ID, "lifecycle"))
		if loadErr == nil {
			workflows = append(workflows, w)
		}
	}
	sort.Slice(workflows, func(i, j int) bool { return workflows[i].ID < workflows[j].ID })
	return Snapshot{Records: records, Audits: audits, Attachments: attachments, Workflows: workflows}, nil
}

func (s Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	err := json.Unmarshal(data, &snapshot)
	return snapshot, err
}

func (s *Store) ImportSnapshot(snapshot Snapshot) error {
	for _, r := range snapshot.Records {
		if err := s.SaveRecord(r); err != nil {
			return err
		}
	}
	for _, e := range snapshot.Audits {
		if err := s.SaveAudit(e); err != nil {
			return err
		}
	}
	for _, w := range snapshot.Workflows {
		if err := s.SaveWorkflow(w); err != nil {
			return err
		}
	}
	for _, a := range snapshot.Attachments {
		if err := s.SaveAttachment(a); err != nil {
			return err
		}
	}
	return nil
}
