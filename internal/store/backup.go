package store

import (
	"encoding/json"
	"fmt"
	"sort"

	"telemetry.local/drone/internal/domain"
)

type BackupManifest struct {
	Format      string `json:"format"`
	RecordCount int    `json:"record_count"`
	AuditCount  int    `json:"audit_count"`
	Checksum    string `json:"checksum"`
}

func (s *Store) Backup() ([]byte, BackupManifest, error) {
	snapshot, err := s.ExportSnapshot()
	if err != nil {
		return nil, BackupManifest{}, err
	}
	payload, err := snapshot.Marshal()
	if err != nil {
		return nil, BackupManifest{}, err
	}
	manifest := BackupManifest{Format: "telemetry-bbolt-v1", RecordCount: len(snapshot.Records), AuditCount: len(snapshot.Audits), Checksum: domain.StableID(string(payload))}
	packet := struct {
		Manifest BackupManifest `json:"manifest"`
		Snapshot Snapshot       `json:"snapshot"`
	}{Manifest: manifest, Snapshot: snapshot}
	data, err := json.Marshal(packet)
	return data, manifest, err
}

func (s *Store) Restore(data []byte) (BackupManifest, error) {
	var packet struct {
		Manifest BackupManifest `json:"manifest"`
		Snapshot Snapshot       `json:"snapshot"`
	}
	if err := json.Unmarshal(data, &packet); err != nil {
		return BackupManifest{}, err
	}
	if packet.Manifest.Format != "telemetry-bbolt-v1" {
		return BackupManifest{}, fmt.Errorf("unsupported backup format")
	}
	if packet.Manifest.RecordCount != len(packet.Snapshot.Records) || packet.Manifest.AuditCount != len(packet.Snapshot.Audits) {
		return BackupManifest{}, fmt.Errorf("backup manifest counts do not match")
	}
	if err := s.ImportSnapshot(packet.Snapshot); err != nil {
		return BackupManifest{}, err
	}
	return packet.Manifest, nil
}

func SortManifestRecords(snapshot Snapshot) {
	sort.Slice(snapshot.Records, func(i, j int) bool { return snapshot.Records[i].ID < snapshot.Records[j].ID })
}
