package unit

import (
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"time"
)

const SimpleDateFormat = "2006-01-02T15:04:05"

func createTicket(status string, createTime time.Time, changelogEntries ...domain.ChangelogEntry) domain.Ticket {
	return domain.Ticket{
		Id:               "Ticket-532",
		Key:              "Ticket-532",
		Title:            "Random title",
		ChangelogEntries: changelogEntries,
		UpdateTime:       domain.BeginingOfTime,
		CreateTime:       createTime,
		State:            status,
		DevStartDate:     -1,
		DevEndDate:       -1,
	}
}

func createChangeLogEntries(entries ...domain.ChangelogEntry) []domain.ChangelogEntry {
	return entries
}

func createChangelogEntry(author string, timestamp time.Time, changes ...domain.Change) domain.ChangelogEntry {
	return domain.CreateChangelogEntry(
		"random", author, timestamp, changes...,
	)
}

func simpleChangelogEntry(field string, prevValue string, newValue string, timestamp time.Time) domain.ChangelogEntry {
	return domain.CreateChangelogEntry(
		"random", "randomAuthor", timestamp, domain.CreateChange(field, prevValue, newValue),
	)
}

func statusChangelogEntry(prevValue string, newValue string, timestamp time.Time) domain.ChangelogEntry {
	return simpleChangelogEntry("status", prevValue, newValue, timestamp)
}

func dirtyDate(dateString string) time.Time {
	date, err := time.Parse(SimpleDateFormat, dateString)
	if err != nil {
		panic(err)
	}

	return date
}
