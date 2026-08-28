package store

import (
	"errors"

	"go.etcd.io/bbolt"
	"telemetry.local/drone/internal/domain"
)

type Transaction struct {
	store   *Store
	records []domain.Record
	audits  []domain.AuditEvent
}

func (s *Store) Begin() *Transaction {
	return &Transaction{store: s, records: make([]domain.Record, 0), audits: make([]domain.AuditEvent, 0)}
}

func (t *Transaction) AddRecord(record domain.Record) error {
	if err := domain.ValidateRecord(record); err != nil {
		return err
	}
	t.records = append(t.records, record)
	return nil
}

func (t *Transaction) AddAudit(event domain.AuditEvent) {
	t.audits = append(t.audits, event)
}

func (t *Transaction) Commit() error {
	if t.store == nil || t.store.db == nil {
		return errors.New("transaction has no store")
	}
	return t.store.db.Update(func(tx *bbolt.Tx) error {
		for _, record := range t.records {
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := tx.Bucket(bucketRecords).Put([]byte(record.ID), data); err != nil {
				return err
			}
		}
		for _, event := range t.audits {
			data, err := encode(event)
			if err != nil {
				return err
			}
			if err := tx.Bucket(bucketAudits).Put([]byte(event.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}
