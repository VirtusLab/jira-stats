package domain

import (
	"sort"
)

func MakeIntervals(ticket Ticket, transitions ...Transition) []TransitionInterval {
	// make sure transition points are sorted in asc order (timestamp wise)
	sort.Slice(transitions, func(i, j int) bool {
		return transitions[i].Timestamp.Before(transitions[j].Timestamp)
	})

	var intervals []TransitionInterval
	startTime := ticket.CreateTime

	for _, t := range transitions {
		currentTransition := TransitionInterval{
			Start:  startTime,
			End:    t.Timestamp,
			State:  t.FromState,
			Author: t.Author,
		}

		intervals = append(intervals, currentTransition)

		startTime = t.Timestamp
	}

	intervals = append(intervals, TransitionInterval{
		Start: startTime,
		End:   END_OF_TIME,
		State: ticket.State,
	})

	return intervals
}
