package review

import (
	"fmt"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Decision struct {
	Allowed bool
	Reason  string
}

func Check(record domain.Record, action string) Decision {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "approve":
		if record.Status == domain.StatusPending {
			return Decision{Allowed: true, Reason: "pending record may be reviewed"}
		}
	case "publish":
		if record.Status == domain.StatusApproved {
			return Decision{Allowed: true, Reason: "approved record may be published"}
		}
	case "archive":
		if record.Status == domain.StatusApproved || record.Status == domain.StatusPublished {
			return Decision{Allowed: true, Reason: "reviewed record may be archived"}
		}
	}
	return Decision{Allowed: false, Reason: fmt.Sprintf("action %s is not permitted for %s", action, record.Status)}
}

func CompareVersions(before, after domain.Record) error {
	if after.Version < before.Version {
		return fmt.Errorf("version regressed from %d to %d", before.Version, after.Version)
	}
	if after.ID != before.ID {
		return fmt.Errorf("record identity changed")
	}
	return nil
}

func ActionLabel(action string) string {
	return strings.ToUpper(strings.TrimSpace(action))
}
