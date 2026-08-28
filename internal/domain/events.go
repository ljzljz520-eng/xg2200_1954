package domain

type EventBuilder struct {
	RecordID string
	Actor    string
	Next     int
}

func NewEventBuilder(recordID, actor string, start int) EventBuilder {
	if start < 1 {
		start = 1
	}
	return EventBuilder{RecordID: recordID, Actor: actor, Next: start}
}

func (b *EventBuilder) Build(action, detail string) AuditEvent {
	e := AuditEvent{ID: AuditID(b.RecordID, b.Next), RecordID: b.RecordID, Action: action, Actor: b.Actor, Detail: detail, Sequence: b.Next}
	b.Next++
	return e
}

func WorkflowFor(recordID, name string, step int, state string) Workflow {
	return Workflow{ID: WorkflowID(recordID, name), RecordID: recordID, Name: name, CurrentStep: step, State: state}
}
