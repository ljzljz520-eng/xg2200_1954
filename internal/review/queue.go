package review

import (
	"sort"

	"telemetry.local/drone/internal/domain"
)

type Queue struct {
	items []domain.Record
}

func NewQueue(records []domain.Record) *Queue {
	items := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if record.Status == domain.StatusPending {
			items = append(items, record)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt < items[j].CreatedAt || (items[i].CreatedAt == items[j].CreatedAt && items[i].ID < items[j].ID)
	})
	return &Queue{items: items}
}

func (q *Queue) Next() (domain.Record, bool) {
	if len(q.items) == 0 {
		return domain.Record{}, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Len() int { return len(q.items) }

func (q *Queue) IDs() []string {
	ids := make([]string, 0, len(q.items))
	for _, item := range q.items {
		ids = append(ids, item.ID)
	}
	return ids
}
