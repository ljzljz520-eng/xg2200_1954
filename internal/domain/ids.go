package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func RecordID(batch string, row int) string {
	return StableID("record", batch, fmt.Sprint(row))
}

func AttachmentID(recordID string, name string) string {
	return StableID("attachment", recordID, name)
}

func AuditID(recordID string, sequence int) string {
	return StableID("audit", recordID, fmt.Sprint(sequence))
}

func WorkflowID(recordID string, name string) string {
	return StableID("workflow", recordID, name)
}

func StableID(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte{0})
		h.Write([]byte(part))
	}
	return hex.EncodeToString(h.Sum(nil))[:20]
}
