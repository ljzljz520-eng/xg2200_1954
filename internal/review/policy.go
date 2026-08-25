package review

import (
	"fmt"
	"strings"

	"telemetry.local/drone/internal/domain"
)

type Gate struct {
	RequireActor  bool
	AllowedActors map[string]bool
}

func DefaultGate() Gate {
	return Gate{RequireActor: true, AllowedActors: map[string]bool{"system-reviewer": true, "reviewer-a": true, "reviewer-b": true, "qa": true, "api": true}}
}

func (g Gate) Check(actor string, record domain.Record, action string) error {
	actor = strings.TrimSpace(actor)
	if g.RequireActor && actor == "" {
		return fmt.Errorf("review actor is required")
	}
	if len(g.AllowedActors) > 0 && !g.AllowedActors[actor] {
		return fmt.Errorf("actor %s is not authorized", actor)
	}
	decision := Check(record, action)
	if !decision.Allowed {
		return fmt.Errorf("review denied: %s", decision.Reason)
	}
	return nil
}

func RequireReason(reason string) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("review reason is required")
	}
	if len(reason) > 300 {
		return fmt.Errorf("review reason too long")
	}
	return nil
}
