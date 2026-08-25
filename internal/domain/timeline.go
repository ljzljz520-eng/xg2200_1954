package domain

import "sort"

type Timeline struct {
	RecordID string
	Events   []AuditEvent
}

func NewTimeline(recordID string, events []AuditEvent) Timeline {
	copyEvents := append([]AuditEvent(nil), events...)
	sort.Slice(copyEvents, func(i, j int) bool { return copyEvents[i].Sequence < copyEvents[j].Sequence })
	return Timeline{RecordID: recordID, Events: copyEvents}
}

func (t Timeline) LastAction() string {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[len(t.Events)-1].Action
}

func (t Timeline) HasAction(action string) bool {
	for _, event := range t.Events {
		if event.Action == action {
			return true
		}
	}
	return false
}

func (t Timeline) SequencesValid() bool {
	for i, event := range t.Events {
		if event.Sequence != i+1 {
			return false
		}
	}
	return true
}

func (t Timeline) Actions() []string {
	actions := make([]string, 0, len(t.Events))
	for _, event := range t.Events {
		actions = append(actions, event.Action)
	}
	return actions
}
