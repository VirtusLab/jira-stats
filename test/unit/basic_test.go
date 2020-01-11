package unit

import (
	jiraProcessor "github.com/VirtusLab/jira-stats/analyzer"
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
	"testing"
	"time"
)

// History with no state transition
func TestCorrectAssignmentForNoTransition(t *testing.T) {
	issue := createJiraIssue(
		changeLog(
			[]jira.ChangelogHistory{
				changeLogHistoryItem(
					"2006-01-01T15:04:05.000-0700",
					[]jira.ChangelogItems{
						changeLogItem("Description", "bam bam", "different"),
					},
				),
			},
		),
	)

	tickets, err := jiraProcessor.BuildModel([]jira.Issue{issue})
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}

	assert.Equal(t, len(tickets), 1, "There should be one analyzer ticket created")
	assert.Equal(t, tickets[0].DevStartDate, domain.EndOfTime.Unix(), "Start date should not be set")
	assert.Equal(t, tickets[0].DevEndDate, domain.BeginingOfTime.Unix(), "End date should not be set")
}

// History array is empty
func TestCorrectAssignmentForNoHistory(t *testing.T) {
	issue := createJiraIssue(
		changeLog(
			[]jira.ChangelogHistory{},
		),
	)

	tickets, err := jiraProcessor.BuildModel([]jira.Issue{issue})
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}

	assert.Equal(t, len(tickets), 1, "There should be one analyzer ticket created")
	assert.Equal(t, tickets[0].DevStartDate, domain.EndOfTime.Unix(), "Start date should not be set")
	assert.Equal(t, tickets[0].DevEndDate, domain.BeginingOfTime.Unix(), "End date should not be set")
}

// simple history with one state transition
func TestCorrectAssignmentForSimpleHistory(t *testing.T) {
	issue := createJiraIssue(
		changeLog(
			[]jira.ChangelogHistory{
				changeLogHistoryItem(
					"2006-01-01T15:04:05.000-0700",
					[]jira.ChangelogItems{
						changeLogItem("Description", "PrevPrev", "Prev"),
					},
				),
				changeLogHistoryItem(
					"2006-01-02T15:04:05.000-0700",
					[]jira.ChangelogItems{
						changeLogItem("Description", "Prev", "Current"),
						changeLogItem("Status", "To Do", "In Development"),
					},
				),
				changeLogHistoryItem(
					"2006-02-02T15:04:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "Ready For Testing")},
				),
			},
		),
	)

	tickets, err := jiraProcessor.BuildModel([]jira.Issue{issue})
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}

	startDate, _ := time.Parse(domain.JiraTimestampFormat, "2006-01-02T15:04:05.000-0700")
	endDate, _ := time.Parse(domain.JiraTimestampFormat, "2006-02-02T15:04:05.000-0700")

	assert.Equal(t, tickets[0].DevStartDate, startDate.Unix(), "Start date should be set correctly")
	assert.Equal(t, tickets[0].DevEndDate, endDate.Unix(), "End date should be set correctly")
}

// two "In Dev" state transitions and some other non-dev related, also order is mixed
func TestCorrectAssignmentForComplicatedHistory(t *testing.T) {
	issue := createJiraIssue(
		changeLog(
			[]jira.ChangelogHistory{
				changeLogHistoryItem(
					"2006-01-02T15:04:05.000-0700",
					[]jira.ChangelogItems{
						changeLogItem("Description", "Prev", "Current"),
						changeLogItem("Status", "To Do", "In Development"),
					},
				),
				changeLogHistoryItem(
					"2006-02-02T15:04:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "Ready For Testing")},
				),
				changeLogHistoryItem(
					"2005-06-03T15:04:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "To Do", "In Development")}, // this happens before transitions above
				),
				changeLogHistoryItem(
					"2005-09-15T15:04:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "To Do")}, // this happens before transitions above
				),
			},
		),
	)

	tickets, err := jiraProcessor.BuildModel([]jira.Issue{issue})
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}

	startDate, _ := time.Parse(domain.JiraTimestampFormat, "2005-06-03T15:04:05.000-0700")
	endDate, _ := time.Parse(domain.JiraTimestampFormat, "2006-02-02T15:04:05.000-0700")

	assert.Equal(t, tickets[0].DevStartDate, startDate.Unix(), "Start date should be set correctly")
	assert.Equal(t, tickets[0].DevEndDate, endDate.Unix(), "End date should be set correctly")
}

func createJiraIssue(changeLog jira.Changelog) jira.Issue {
	return jira.Issue{
		ID:  "11232",
		Key: "ABC-112",
		Fields: &jira.IssueFields{
			Summary: "Ticket summary",
			Status: &jira.Status{
				Name: "To Do",
			},
		},
		Changelog: &changeLog,
	}
}

func changeLog(changeLogHistories []jira.ChangelogHistory) jira.Changelog {
	return jira.Changelog{
		Histories: changeLogHistories,
	}
}

func changeLogHistoryItem(created string, changeLogItems []jira.ChangelogItems) jira.ChangelogHistory {
	return jira.ChangelogHistory{
		Author:  jira.User{Name: "test.user"},
		Created: created,
		Items:   changeLogItems,
	}
}

func changeLogItem(field string, from string, to string) jira.ChangelogItems {
	return jira.ChangelogItems{
		Field:      field,
		FromString: from,
		ToString:   to,
	}
}
