package domain

import (
	"fmt"
	"strings"
)

type Schema struct {
	Name    string
	Version int
	Fields  []string
}

func RecordSchema() Schema {
	return Schema{Name: "telemetry-record", Version: 1, Fields: []string{"id", "title", "status", "payload", "version"}}
}

func (s Schema) ValidateFields(values map[string]string) []string {
	issues := make([]string, 0)
	for _, field := range s.Fields {
		if strings.TrimSpace(values[field]) == "" {
			issues = append(issues, fmt.Sprintf("missing-%s", field))
		}
	}
	return issues
}

func (s Schema) Compatible(version int) bool {
	return version == s.Version
}

func (s Schema) Describe() string {
	return fmt.Sprintf("%s v%d (%d fields)", s.Name, s.Version, len(s.Fields))
}
