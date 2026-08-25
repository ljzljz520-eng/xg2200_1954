package store

import (
	"encoding/json"
	"errors"
	"sort"

	"go.etcd.io/bbolt"
	"telemetry.local/drone/internal/domain"
)

var bucketRecords = []byte("records")
var bucketAudits = []byte("audits")
var bucketWorkflows = []byte("workflows")
var bucketAttachments = []byte("attachments")

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{bucketRecords, bucketAudits, bucketWorkflows, bucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func encode(value any) ([]byte, error) {
	return json.Marshal(value)
}

func decode(data []byte, value any) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}
	return json.Unmarshal(data, value)
}

func put(bucket []byte, key string, value any) func(*bbolt.Tx) error {
	return func(tx *bbolt.Tx) error {
		data, err := encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket(bucket).Put([]byte(key), data)
	}
}

func (s *Store) SaveRecord(r domain.Record) error {
	if err := domain.ValidateRecord(r); err != nil {
		return err
	}
	return s.db.Update(put(bucketRecords, r.ID, r))
}

func (s *Store) LoadRecord(id string) (domain.Record, error) {
	var r domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucketRecords).Get([]byte(id))
		if data == nil {
			return errors.New("record not found")
		}
		return decode(data, &r)
	})
	return r, err
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, data []byte) error {
			var r domain.Record
			if err := decode(data, &r); err != nil {
				return err
			}
			result = append(result, r)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) SaveAudit(e domain.AuditEvent) error {
	data, err := encode(e)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketAudits).Put([]byte(e.ID), data) })
}

func (s *Store) ListAudits(recordID string) ([]domain.AuditEvent, error) {
	items := make([]domain.AuditEvent, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAudits).ForEach(func(_, data []byte) error {
			var e domain.AuditEvent
			if err := decode(data, &e); err != nil {
				return err
			}
			if recordID == "" || e.RecordID == recordID {
				items = append(items, e)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].Sequence < items[j].Sequence })
	return items, err
}

func (s *Store) SaveWorkflow(w domain.Workflow) error {
	return s.db.Update(put(bucketWorkflows, w.ID, w))
}

func (s *Store) LoadWorkflow(id string) (domain.Workflow, error) {
	var w domain.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucketWorkflows).Get([]byte(id))
		if data == nil {
			return errors.New("workflow not found")
		}
		return decode(data, &w)
	})
	return w, err
}

func (s *Store) SaveAttachment(a domain.Attachment) error {
	return s.db.Update(put(bucketAttachments, a.ID, a))
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	items := make([]domain.Attachment, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAttachments).ForEach(func(_, data []byte) error {
			var a domain.Attachment
			if err := decode(data, &a); err != nil {
				return err
			}
			if recordID == "" || a.RecordID == recordID {
				items = append(items, a)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, err
}

func (s *Store) Search(f domain.Filter) ([]domain.Record, error) {
	all, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Record, 0, len(all))
	for _, r := range all {
		if domain.Match(r, f) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
