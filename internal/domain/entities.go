package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusApproved  Status = "approved"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Record struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    Status `json:"status"`
	Payload   string `json:"payload"`
	Version   int    `json:"version"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	Sequence int    `json:"sequence"`
}

type Workflow struct {
	ID          string `json:"id"`
	RecordID    string `json:"record_id"`
	Name        string `json:"name"`
	CurrentStep int    `json:"current_step"`
	State       string `json:"state"`
	StartedAt   int64  `json:"started_at"`
	CompletedAt int64  `json:"completed_at"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Digest   string `json:"digest"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("record title is required")
	}
	for _, ch := range r.Title {
		if ch < 32 || ch == 127 {
			return errors.New("record title contains control characters")
		}
	}
	if strings.TrimSpace(r.Payload) == "" {
		return errors.New("record payload is required")
	}
	if len(r.Payload) > 8192 {
		return errors.New("record payload is too large")
	}
	if r.Version < 1 {
		return errors.New("record version must be positive")
	}
	if r.Status == "" {
		return errors.New("record status is required")
	}
	return nil
}

func Transition(from, to Status) error {
	if from == StatusPending && to == StatusApproved {
		return nil
	}
	if from == StatusApproved && to == StatusPublished {
		return nil
	}
	if from == StatusApproved && to == StatusArchived {
		return nil
	}
	if from == StatusPublished && to == StatusArchived {
		return nil
	}
	return fmt.Errorf("invalid transition %s to %s", from, to)
}

func IsActive(r Record) bool {
	return r.Status == StatusApproved || r.Status == StatusPublished
}

func NormalizeTitle(title string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(title)), " ")
}

func AllowedTitle(title string) bool {
	normalized := NormalizeTitle(title)
	if normalized == "" {
		return false
	}
	if strings.Contains(strings.ToLower(normalized), "expired") {
		return false
	}
	return len(normalized) <= 160
}
