package domain

import (
	"fmt"
	"strings"
)

type Policy struct {
	MaxPayload   int
	MaxTitle     int
	RequireID    bool
	AllowedKinds map[string]bool
}

func DefaultPolicy() Policy {
	return Policy{MaxPayload: 8192, MaxTitle: 160, RequireID: true, AllowedKinds: map[string]bool{"telemetry": true, "signature": true, "manifest": true}}
}

func (p Policy) CheckRecord(record Record) error {
	if p.RequireID && strings.TrimSpace(record.ID) == "" {
		return fmt.Errorf("policy requires an id")
	}
	if len(record.Title) > p.MaxTitle {
		return fmt.Errorf("title exceeds policy limit")
	}
	if len(record.Payload) > p.MaxPayload {
		return fmt.Errorf("payload exceeds policy limit")
	}
	return ValidateRecord(record)
}

func (p Policy) CheckAttachment(attachment Attachment) error {
	if strings.TrimSpace(attachment.RecordID) == "" || strings.TrimSpace(attachment.Name) == "" {
		return fmt.Errorf("attachment identity is required")
	}
	if attachment.Size < 0 {
		return fmt.Errorf("attachment size cannot be negative")
	}
	if len(p.AllowedKinds) > 0 && !p.AllowedKinds[attachment.Kind] {
		return fmt.Errorf("attachment kind %s is not allowed", attachment.Kind)
	}
	return nil
}

func (p Policy) Explain(record Record) []string {
	issues := make([]string, 0)
	if strings.TrimSpace(record.ID) == "" {
		issues = append(issues, "missing-id")
	}
	if strings.TrimSpace(record.Title) == "" {
		issues = append(issues, "missing-title")
	}
	if strings.Contains(strings.ToLower(record.Title), "expired") {
		issues = append(issues, "expired-title")
	}
	if len(record.Payload) > p.MaxPayload {
		issues = append(issues, "payload-too-large")
	}
	return issues
}
