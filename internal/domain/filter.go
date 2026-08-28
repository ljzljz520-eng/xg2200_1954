package domain

import "strings"

type Filter struct {
	Query  string
	Status Status
}

func Match(r Record, f Filter) bool {
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Query == "" {
		return true
	}
	q := strings.ToLower(strings.TrimSpace(f.Query))
	return strings.Contains(strings.ToLower(r.ID), q) || strings.Contains(strings.ToLower(r.Title), q) || strings.Contains(strings.ToLower(r.Payload), q)
}

func SortRecords(records []Record) {
	for i := 1; i < len(records); i++ {
		current := records[i]
		j := i - 1
		for j >= 0 && records[j].ID > current.ID {
			records[j+1] = records[j]
			j--
		}
		records[j+1] = current
	}
}
