package unit

import (
	"github.com/stretchr/testify/assert"
	"jira-stats/jira-stats-serverless/analyzer/domain"
	"testing"
)

func TestSimpleTransition(t *testing.T) {
	ticket := createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	transitions := domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-01-02T00:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-02-05T23:59:59")),
	)
	assert.Equal(t, len(transitions), 3, "Number of generated intervals incorrect")

	assert.Equal(t, transitions[0], domain.TransitionInterval{
		Start: dirtyDate("2018-01-01T00:00:00"),
		End:   dirtyDate("2018-01-02T00:00:00"),
		State: "To Do",
	}, "Incorrect number of dev days calculated")

	assert.Equal(t, transitions[1], domain.TransitionInterval{
		Start: dirtyDate("2018-01-02T00:00:00"),
		End:   dirtyDate("2018-02-05T23:59:59"),
		State: "In Development",
	}, "Incorrect number of dev days calculated")

	assert.Equal(t, transitions[2], domain.TransitionInterval{
		Start: dirtyDate("2018-02-05T23:59:59"),
		End:   domain.END_OF_TIME,
		State: "In Review",
	}, "Incorrect number of dev days calculated")
}

func TestEmptyTransition(t *testing.T) {
	ticket := createTicket("Open", dirtyDate("2018-01-01T00:00:00"))
	transitions := domain.MakeIntervals(ticket)
	assert.Equal(t, len(transitions), 1, "Number of generated intervals incorrect")

	assert.Equal(t, transitions[0], domain.TransitionInterval{
		Start: dirtyDate("2018-01-01T00:00:00"),
		End:   domain.END_OF_TIME,
		State: "Open",
	}, "Incorrect number of dev days calculated")
}
