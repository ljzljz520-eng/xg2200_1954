package importer

import "strings"

func DeterministicBatch() string {
	return strings.Join([]string{
		"alpha|Current Flight Plan|cipher-a|telemetry",
		"bravo|Night Survey Route|cipher-b|telemetry",
		"charlie|expired flight plan|cipher-c|telemetry",
	}, "\n")
}

func SecondBatch() string {
	return "delta|Expired Flight Plan|cipher-d|telemetry"
}
