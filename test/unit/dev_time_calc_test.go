package unit

import (
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test with single day work
func TestOneDayInterval(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-01-02T13:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-01-02T13:15:59")),
	)

	calculator := domain.DaysCalculator{}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.25, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-01-01T00:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-01-30T19:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-01-30T22:15:59")),
	)
	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.5, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-03-01T00:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-03-30T19:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-03-31T00:10:59")),
	)
	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.75, days, "Incorrect number of dev hours calculated")
}

// Test times completely outside of given boundaries
func TestIntervalsOutsideOfBoundaries(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2017-10-31T19:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2017-11-30T21:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2017-12-31T19:00:00")),
	)
	calculator := domain.DaysCalculator{}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.0, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2019-03-31T22:15:59"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2019-03-31T22:23:59")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2019-04-01T13:30:13")),
	)

	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.0, days, "Incorrect number of dev hours calculated")
}

// Test times slashed by limits of given time boundaries
func TestOneDayIntervalLimitedByBounds(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2017-12-31T19:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2017-12-31T21:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-01-01T01:30:13")),
	)

	calculator := domain.DaysCalculator{}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.25, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-03-31T22:15:59"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-03-31T21:23:59")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-04-01T13:30:13")),
	)

	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.5, days, "Incorrect number of dev hours calculated")
}

// Tests simple scenarios over several days
func TestSimpleScenario(t *testing.T) {
	startDate := dirtyDate("2018-01-01T00:00:00")
	endDate := dirtyDate("2018-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2018-02-01T09:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-02-01T10:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-02-02T19:00:00")),
	)

	calculator := domain.DaysCalculator{}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 2.0, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2018-02-01T09:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2018-02-05T13:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2018-02-06T19:00:00")),
	)

	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.5, days, "Incorrect number of dev hours calculated")
}

// Tests calculation only of working days
func TestSkippingWeekendDay(t *testing.T) {
	startDate := dirtyDate("2020-01-01T00:00:00")
	endDate := dirtyDate("2020-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2020-02-01T09:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2020-02-02T09:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2020-02-03T19:00:00")),
	)
	calculator := domain.DaysCalculator{
		ClockNow: func() time.Time {
			return dirtyDate("2020-12-31T00:00:00")
		},
	}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.0, days, "Incorrect number of dev hours calculated")

	ticket = createTicket("In Review", dirtyDate("2020-01-01T09:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2020-01-13T11:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2020-01-29T09:00:00")),
	)

	days = calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 12+0.5, days, "Incorrect number of dev hours calculated")
}

// Tests calculation multiple dev intervals
func TestMultipleDevIntervals(t *testing.T) {
	startDate := dirtyDate("2020-01-01T00:00:00")
	endDate := dirtyDate("2020-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2020-02-01T09:00:00"))
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2020-02-02T09:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2020-02-03T19:00:00")),
		statusChangelogEntry("In Review", "In Development", dirtyDate("2020-02-04T07:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2020-02-06T11:00:00")),
	)

	calculator := domain.DaysCalculator{
		ClockNow: func() time.Time {
			return dirtyDate("2020-12-31T00:00:00")
		},
	}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 1.0+2.5, days, "Incorrect number of dev hours calculated")
}

// Tests skipping calculations
func TestSkippingTickets(t *testing.T) {
	startDate := dirtyDate("2020-01-01T00:00:00")
	endDate := dirtyDate("2020-03-31T23:59:59")

	ticket := createTicket("In Review", dirtyDate("2020-02-01T09:00:00"))
	ticket.Type = "Epic"
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", dirtyDate("2020-02-02T09:00:00")),
		statusChangelogEntry("In Development", "In Review", dirtyDate("2020-02-03T19:00:00")),
	)

	calculator := domain.DaysCalculator{}
	days := calculator.CalculateDevDays(ticket, startDate, endDate)
	assert.Equal(t, 0.0, days, "Epic should not be taken into account")
}

// Tests interval being cut off by current date
func TestLimitedByCurrentDate(t *testing.T) {
	currentTime := time.Date(2020, 1, 10, 8, 0, 0, 0, time.UTC)
	createTime := currentTime.AddDate(0, 0, -5)
	devStartTime := currentTime.AddDate(0, 0, -3)
	intervalStart := currentTime.AddDate(0, 0, -10)
	intervalEnd := currentTime.AddDate(0, 0, 10)

	ticket := createTicket("In Development", createTime)
	ticket.ChangelogEntries = createChangeLogEntries(
		statusChangelogEntry("To Do", "In Development", devStartTime),
	)

	calculator := domain.DaysCalculator{
		ClockNow: func() time.Time {
			return currentTime
		},
	}
	days := calculator.CalculateDevDays(ticket, intervalStart, intervalEnd)
	assert.Equal(t, 3.5, days, "Hours should not be calculated beyond current date")
}
