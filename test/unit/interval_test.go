package unit

import (
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleTransition(t *testing.T) {
	ticket := createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-01-02T00:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-02-05T23:59:59")),
	)
	transitions := domain.MakeIntervals(ticket)

	assert.Equal(t, len(transitions), 3, "Number of generated intervals incorrect")

	assert.Equal(t, domain.TransitionInterval{
		Start: dirtyDate("2018-01-01T00:00:00"),
		End:   dirtyDate("2018-01-02T00:00:00"),
		State: "To Do",
	}, transitions[0], "Incorrect number of dev days calculated")

	assert.Equal(t, domain.TransitionInterval{
		Start: dirtyDate("2018-01-02T00:00:00"),
		End:   dirtyDate("2018-02-05T23:59:59"),
		State: "In Development",
	}, transitions[1], "Incorrect number of dev days calculated")

	assert.Equal(t, domain.TransitionInterval{
		Start: dirtyDate("2018-02-05T23:59:59"),
		End:   domain.EndOfTime,
		State: "In Review",
	}, transitions[2], "Incorrect number of dev days calculated")
}

func TestEmptyTransition(t *testing.T) {
	ticket := createTicket("Open", dirtyDate("2018-01-01T00:00:00"))
	transitions := domain.MakeIntervals(ticket)
	assert.Equal(t, len(transitions), 1, "Number of generated intervals incorrect")

	assert.Equal(t, domain.TransitionInterval{
		Start: dirtyDate("2018-01-01T00:00:00"),
		End:   domain.EndOfTime,
		State: "Open",
	}, transitions[0], "Incorrect number of dev days calculated")
}
