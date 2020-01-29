package analyzer

import (
	"encoding/json"
	"log"
	"time"
)

type TimingEvent struct {
	EventType    string            `json:"type"`
	OpName       string            `json:"opName"`
	TimeInMillis int64             `json:"time"`
	Params       map[string]string `json:"params"`
}

func TimeTrack(start time.Time, name string) {
	timeTrackParams(start, name, map[string]string{})
}

func timeTrackParams(start time.Time, name string, params map[string]string) {
	elapsed := time.Since(start).Milliseconds()
	bytes, err := json.Marshal(TimingEvent{
		EventType:    "timing",
		OpName:       name,
		TimeInMillis: elapsed,
		Params:       params,
	})

	if err != nil {
		panic("Unexpected marshalling error")
	}

	log.Println(string(bytes))
}
