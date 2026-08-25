package report

import (
	"sort"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type AuditView struct {
	Sequence int    `json:"sequence"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
}

func BuildAuditView(events []domain.AuditEvent) []AuditView {
	copyEvents := append([]domain.AuditEvent(nil), events...)
	sort.Slice(copyEvents, func(i, j int) bool { return copyEvents[i].Sequence < copyEvents[j].Sequence })
	views := make([]AuditView, 0, len(copyEvents))
	for _, event := range copyEvents {
		views = append(views, AuditView{Sequence: event.Sequence, Action: event.Action, Actor: event.Actor, Detail: event.Detail})
	}
	return views
}

func FilterActions(events []domain.AuditEvent, action string) []domain.AuditEvent {
	filtered := make([]domain.AuditEvent, 0)
	for _, event := range events {
		if strings.EqualFold(event.Action, action) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func CountActions(events []domain.AuditEvent) map[string]int {
	counts := make(map[string]int)
	for _, event := range events {
		counts[event.Action]++
	}
	return counts
}
