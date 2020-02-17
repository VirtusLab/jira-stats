package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Transition struct {
	prevState string
	newState  string
	timestamp time.Time
}

type TransitionInterval struct {
	Start time.Time
	End   time.Time
	State string
}

func (t *TransitionInterval) ToString() string {
	return fmt.Sprintf("TransitionInterval [Start: %s, End: %s, State: %s]",
		t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339), t.State)
}

func MakeIntervals(ticket Ticket) []TransitionInterval {
	transitions := extractStateChanges(ticket.ChangelogEntries)

	var intervals []TransitionInterval
	startTime := time.Unix(ticket.CreateTime, 0).UTC()

	for _, t := range transitions {
		currentTransition := TransitionInterval{
			Start: startTime,
			End:   t.timestamp,
			State: t.prevState,
		}

		intervals = append(intervals, currentTransition)

		startTime = t.timestamp
	}

	intervals = append(intervals, TransitionInterval{
		Start: startTime,
		End:   EndOfTime,
		State: ticket.State,
	})

	return intervals
}

func extractStateChanges(changeLogEntries []ChangelogEntry) []Transition {
	transitions := make([]Transition, 0)

	for _, entry := range changeLogEntries {
		change, err := getStatusChange(entry)
		if err == nil {
			transitions = append(transitions, Transition{change.From, change.To, entry.Created})
		}
	}

	// make sure transition points are sorted in asc order (timestamp wise)
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].timestamp.Before(transitions[j].timestamp)
	})

	return transitions
}

func getStatusChange(changeLogEntry ChangelogEntry) (Change, error) {
	for _, change := range changeLogEntry.Changes {
		if strings.ToLower(change.Field) == StatusField {
			return change, nil
		}
	}

	return Change{}, errors.New("No Status transition available in this changelog entry")
}
