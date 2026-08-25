package review

import (
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type History struct {
	RecordID string
	Actions  []string
	Actors   []string
}

func BuildHistory(recordID string, events []domain.AuditEvent) History {
	copyEvents := append([]domain.AuditEvent(nil), events...)
	sort.Slice(copyEvents, func(i, j int) bool { return copyEvents[i].Sequence < copyEvents[j].Sequence })
	history := History{RecordID: recordID, Actions: make([]string, 0, len(copyEvents)), Actors: make([]string, 0, len(copyEvents))}
	for _, event := range copyEvents {
		history.Actions = append(history.Actions, event.Action)
		history.Actors = append(history.Actors, event.Actor)
	}
	return history
}

func (h History) Contains(action string) bool {
	for _, value := range h.Actions {
		if strings.EqualFold(value, action) {
			return true
		}
	}
	return false
}

func (h History) Complete() bool {
	return h.Contains("import") && (h.Contains("approve") || h.Contains("publish")) && h.Contains("archive")
}
