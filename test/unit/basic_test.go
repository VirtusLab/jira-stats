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
	createTime, _ := time.Parse(domain.JiraTimestampFormat, "2006-01-01T13:04:05.000-0700")
	issue := createJiraIssue(
		createTime,
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
	assert.Equal(t, createTime.Unix(), tickets[0].CreateTime, "End date should not be set")
	assert.Equal(t, domain.EndOfTime.Unix(), tickets[0].CloseTime, "End date should not be set")
}

// History array is empty
func TestCorrectAssignmentForNoHistory(t *testing.T) {
	issue := createJiraIssue(
		dirtyDate("2006-01-01T13:04:05"),
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
	assert.Equal(t, domain.EndOfTime.Unix(), tickets[0].CloseTime, "End date should not be set")
}

// simple history with one state transition
func TestCorrectAssignmentForSimpleHistory(t *testing.T) {
	createTime, _ := time.Parse(domain.JiraTimestampFormat, "2006-01-01T13:04:05.000-0700")

	issue := createJiraIssue(
		createTime,
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
						changeLogItem("Status", "To Do", "Closed"),
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

	assert.Equal(t, createTime.Unix(), tickets[0].CreateTime, "Create date should be set correctly")
	assert.Equal(t, domain.EndOfTime.Unix(), tickets[0].CloseTime, "End date should be set correctly")
}

// two "In Dev" state transitions and some other non-dev related, also order is mixed
func TestCorrectAssignmentForComplicatedHistory(t *testing.T) {
	issue := createJiraIssue(
		dirtyDate("2006-01-01T13:04:05"),
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
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "On Live")},
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

	endTime, _ := time.Parse(domain.JiraTimestampFormat, "2006-02-02T15:04:05.000-0700")
	assert.Equal(t, endTime.Unix(), tickets[0].CloseTime, "End date should be set correctly")
}

// two "In Dev" with other state in-between
func TestTwoInDevStates(t *testing.T) {
	firstInDevTransitionIn := "2006-02-02T15:04:05.000-0700"
	lastInDevTransitionOut := "2006-05-12T15:04:05.000-0700"

	issue := createJiraIssue(
		dirtyDate("2006-01-01T13:04:05"),
		changeLog(
			[]jira.ChangelogHistory{
				changeLogHistoryItem(
					"2005-01-02T15:04:05.000-0700",
					[]jira.ChangelogItems{
						changeLogItem("Description", "Prev", "Current"),
					},
				),
				changeLogHistoryItem( // 1st In Dev state
					firstInDevTransitionIn,
					[]jira.ChangelogItems{changeLogItem("Status", "To Do", "In Development")},
				),
				changeLogHistoryItem( // different state
					"2006-03-03T15:04:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "Testing")},
				),
				changeLogHistoryItem( // In Dev again
					"2006-05-04T12:15:05.000-0700",
					[]jira.ChangelogItems{changeLogItem("Status", "Testing", "In Development")},
				),
				changeLogHistoryItem( // and another state again
					lastInDevTransitionOut,
					[]jira.ChangelogItems{changeLogItem("Status", "In Development", "Closed")},
				),
			},
		),
	)

	tickets, err := jiraProcessor.BuildModel([]jira.Issue{issue})
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}

	endDate, _ := time.Parse(domain.JiraTimestampFormat, lastInDevTransitionOut)
	assert.Equal(t, endDate.Unix(), tickets[0].CloseTime, "End date should be set correctly")
}

func createJiraIssue(createTime time.Time, changeLog jira.Changelog) jira.Issue {
	return jira.Issue{
		ID:  "11232",
		Key: "ABC-112",
		Fields: &jira.IssueFields{
			Summary: "Ticket summary",
			Status: &jira.Status{
				Name: "To Do",
			},
			Created: jira.Time(createTime),
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
