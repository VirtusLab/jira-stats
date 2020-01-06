package unit

import (
	"github.com/stretchr/testify/assert"
	"jira-stats/jira-stats-serverless/analyzer/domain"
	"testing"
)

// Test with single day work
func TestOneDayInterval(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-01-02T13:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-01-02T13:15:59")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.25, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-01-30T19:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-01-30T22:15:59")),
	)
	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.5, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-03-01T00:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-03-30T19:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-03-31T00:10:59")),
	)
	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.75, hours, "Incorrect number of dev hours calculated")
}

// Test times completely outside of given boundaries
func TestIntervalsOutsideOfBoundaries(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2017-10-31T19:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2017-11-30T21:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2017-12-31T19:00:00")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.0, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2019-03-31T22:15:59"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2019-03-31T22:23:59")),
		createTransition("In Development", "In Review", dirtyDate("2019-04-01T13:30:13")),
	)

	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.0, hours, "Incorrect number of dev hours calculated")
}

// Test times slashed by limits of given time boundaries
func TestOneDayIntervalLimitedByBounds(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2017-12-31T19:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2017-12-31T21:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-01-01T01:30:13")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.25, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-03-31T22:15:59"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-03-31T21:23:59")),
		createTransition("In Development", "In Review", dirtyDate("2018-04-01T13:30:13")),
	)

	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.5, hours, "Incorrect number of dev hours calculated")
}

// Tests simple scenarios over several days
func TestSimpleScenario(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2018-02-01T09:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-02-01T10:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-02-02T19:00:00")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 2.0, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-02-01T09:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2018-02-05T13:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2018-02-06T19:00:00")),
	)

	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.5, hours, "Incorrect number of dev hours calculated")
}

// Tests calculation only of working days
func TestSkippingWeekendDay(t *testing.T) {
	startDate := dirtyDate("2020-01-01T00:00:00")
	endDate := dirtyDate("2020-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2020-02-01T09:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2020-02-02T09:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2020-02-03T19:00:00")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.0, hours, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2020-01-01T09:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2020-01-13T11:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2020-01-29T09:00:00")),
	)

	hours = domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 12+0.5, hours, "Incorrect number of dev hours calculated")
}

// Tests calculation multiple dev intervals
func TestMultipleDevIntervals(t *testing.T) {
	startDate := dirtyDate("2020-01-01T00:00:00")
	endDate := dirtyDate("2020-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2020-02-01T09:00:00"))
	ticket.Transitions = domain.MakeIntervals(ticket,
		createTransition("To Do", "In Development", dirtyDate("2020-02-02T09:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2020-02-03T19:00:00")),
		createTransition("In Review", "In Development", dirtyDate("2020-02-04T07:00:00")),
		createTransition("In Development", "In Review", dirtyDate("2020-02-06T11:00:00")),
	)

	hours := domain.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.0+2.5, hours, "Incorrect number of dev hours calculated")
}
