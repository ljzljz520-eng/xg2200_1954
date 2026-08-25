package service

import (
	"fmt"
	"sort"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/importer"
)

type ImportPlan struct {
	BatchID  string
	Records  []domain.Record
	Rejected []string
	Manifest importer.Manifest
	DryRun   bool
}

func (s *Service) PlanImport(result importer.Result, batch string, dryRun bool) ImportPlan {
	valid := make([]domain.Record, 0, len(result.Records))
	for _, record := range result.Records {
		if domain.AllowedTitle(record.Title) {
			valid = append(valid, record)
		}
	}
	sort.Slice(valid, func(i, j int) bool { return valid[i].ID < valid[j].ID })
	result.Records = valid
	return ImportPlan{BatchID: batch, Records: valid, Rejected: append([]string(nil), result.Rejected...), Manifest: importer.BuildManifest(batch, result), DryRun: dryRun}
}

func (p ImportPlan) Validate() error {
	if p.BatchID == "" {
		return fmt.Errorf("batch id is required")
	}
	seen := make(map[string]bool)
	for _, record := range p.Records {
		if seen[record.ID] {
			return fmt.Errorf("duplicate record id %s", record.ID)
		}
		seen[record.ID] = true
		if err := domain.ValidateRecord(record); err != nil {
			return err
		}
	}
	return nil
}

func (p ImportPlan) Summary() string {
	mode := "commit"
	if p.DryRun {
		mode = "dry-run"
	}
	return fmt.Sprintf("batch=%s mode=%s accepted=%d rejected=%d digest=%s", p.BatchID, mode, len(p.Records), len(p.Rejected), p.Manifest.Digest)
}
