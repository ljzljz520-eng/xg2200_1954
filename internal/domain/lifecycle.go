package domain

import (
	"fmt"
	"sort"
)

type Lifecycle struct {
	Record    Record
	Workflow  Workflow
	Audits    []AuditEvent
	Completed bool
}

func AssembleLifecycle(record Record, workflow Workflow, audits []AuditEvent) Lifecycle {
	ordered := append([]AuditEvent(nil), audits...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	completed := workflow.State == "archived" || workflow.State == "published"
	return Lifecycle{Record: record, Workflow: workflow, Audits: ordered, Completed: completed}
}

func (l Lifecycle) Validate() error {
	if err := ValidateRecord(l.Record); err != nil {
		return err
	}
	if l.Workflow.RecordID != l.Record.ID {
		return fmt.Errorf("workflow record mismatch")
	}
	if len(l.Audits) == 0 {
		return fmt.Errorf("lifecycle requires audit history")
	}
	for index, event := range l.Audits {
		if event.RecordID != l.Record.ID {
			return fmt.Errorf("audit record mismatch")
		}
		if event.Sequence != index+1 {
			return fmt.Errorf("audit sequence gap")
		}
	}
	return nil
}

func (l Lifecycle) Progress() int {
	if l.Workflow.CurrentStep > 4 {
		return 4
	}
	if l.Workflow.CurrentStep < 0 {
		return 0
	}
	return l.Workflow.CurrentStep
}
