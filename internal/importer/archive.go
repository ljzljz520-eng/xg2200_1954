package importer

import (
	"encoding/json"
	"fmt"
	"sort"

	"telemetry.local/drone/internal/domain"
)

type Archive struct {
	BatchID string
	Rows    []domain.Record
	Files   []domain.Attachment
	Notes   []string
}

func NewArchive(batch string, result Result) Archive {
	rows := append([]domain.Record(nil), result.Records...)
	files := append([]domain.Attachment(nil), result.Attachments...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	sort.Slice(files, func(i, j int) bool { return files[i].ID < files[j].ID })
	notes := append([]string(nil), result.Rejected...)
	return Archive{BatchID: batch, Rows: rows, Files: files, Notes: notes}
}

func (a Archive) Validate() error {
	if a.BatchID == "" {
		return fmt.Errorf("archive batch is required")
	}
	for _, row := range a.Rows {
		if err := domain.ValidateRecord(row); err != nil {
			return err
		}
	}
	for _, file := range a.Files {
		if file.RecordID == "" || file.Digest == "" {
			return fmt.Errorf("archive attachment is incomplete")
		}
	}
	return nil
}

func (a Archive) JSON() ([]byte, error) {
	return json.Marshal(a)
}

func ReadArchive(data []byte) (Archive, error) {
	var archive Archive
	if err := json.Unmarshal(data, &archive); err != nil {
		return Archive{}, err
	}
	if err := archive.Validate(); err != nil {
		return Archive{}, err
	}
	return archive, nil
}
