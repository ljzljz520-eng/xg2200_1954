package domain

type Metrics struct {
	PayloadBytes int
	TitleBytes   int
	Versions     int
	Transitions  int
}

func Measure(record Record, events []AuditEvent) Metrics {
	transitions := 0
	for _, event := range events {
		if event.RecordID == record.ID {
			transitions++
		}
	}
	return Metrics{PayloadBytes: len(record.Payload), TitleBytes: len(record.Title), Versions: record.Version, Transitions: transitions}
}

func Aggregate(records []Record) Metrics {
	result := Metrics{}
	for _, record := range records {
		result.PayloadBytes += len(record.Payload)
		result.TitleBytes += len(record.Title)
		if record.Version > result.Versions {
			result.Versions = record.Version
		}
		if record.Status != StatusPending {
			result.Transitions++
		}
	}
	return result
}
