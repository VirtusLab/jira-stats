package unit

import (
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"time"
)

func createTicket(status string, createTime time.Time, transitions ...domain.TransitionInterval) domain.Ticket {
	return domain.Ticket{
		Id:           "Ticket-532",
		Key:          "Ticket-532",
		Title:        "Random title",
		Transitions:  transitions,
		UpdateTime:   domain.BEGINING_OF_TIME,
		CreateTime:   createTime,
		State:        status,
		DevStartDate: -1,
		DevEndDate:   -1,
	}
}

func createTransition(from string, to string, timestamp time.Time) domain.Transition {

	return domain.Transition{
		FromState: from,
		ToState:   to,
		Timestamp: timestamp,
	}
}

func dirtyDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		panic(err)
	}

	return date
}
