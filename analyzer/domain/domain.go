package domain

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"strings"
	"time"
)

var BEGINING_OF_TIME = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
var END_OF_TIME = time.Date(9999, 12, 31, 23, 59, 59, 999, time.UTC)

const JiraTimestampFormat = "2006-01-02T15:04:05.000-0700"
const JiraUpdateTimestampFormat = "2006-01-02T15:04:05-0700"

const JiraFilterFormat = "2006-01-02 15:04"
const DayFormat = "2006-01-02"

type Ticket struct {
	Id          string
	Key         string
	State       string
	Type        string
	Title       string
	Transitions []TransitionInterval
	UpdateTime  time.Time
	CreateTime  time.Time

	DevStartDate int64
	DevEndDate   int64
}

func (t *Ticket) Project() string {
	dashIdx := strings.LastIndex(t.Key, "-")
	return t.Key[0:dashIdx]
}

type Transition struct {
	FromState string
	ToState   string
	Timestamp time.Time
	Author    string
}

type TransitionInterval struct {
	Start  time.Time
	End    time.Time
	State  string
	Author string
}

func (t *TransitionInterval) ToString() string {
	return fmt.Sprintf("TransitionInterval [Start: %s, End: %s, State: %s, Author: %s]",
		t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339), t.State, t.Author)
}

type ConfigItem struct {
	ConfigName  string
	ConfigValue string
}

type CsvContents struct {
	Header []string
	Rows   []CsvRow
}

type CsvRow struct {
	Entries []string
}

// Prints CSV structure to string
func (csv CsvContents) ToString() string {
	csvList := make([]string, 0)
	csvList = append(csvList, strings.Join(csv.Header, ","))

	for _, row := range csv.Rows {
		csvList = append(csvList, strings.Join(row.Entries, ","))
	}

	return strings.Join(csvList, "\n")
}

func JiraToDomain(jiraIssue jira.Issue) (Ticket, error) {

	transitions := make([]Transition, 0)

	devStartDate := END_OF_TIME
	devEndDate := BEGINING_OF_TIME

	for _, historyItem := range jiraIssue.Changelog.Histories {
		for _, changeItem := range historyItem.Items {
			if strings.ToLower(changeItem.Field) == "status" {

				timestamp, err := time.Parse(JiraTimestampFormat, historyItem.Created)
				if err != nil {
					return Ticket{}, tracerr.Wrap(err)
				}

				transitions = append(transitions, Transition{
					FromState: changeItem.FromString,
					ToState:   changeItem.ToString,
					Timestamp: timestamp,
					Author:    historyItem.Author.Name,
				})

				if changeItem.ToString == "In Development" && devStartDate.After(timestamp) {
					devStartDate = timestamp
				}

				if changeItem.FromString == "In Development" && devEndDate.Before(timestamp) {
					devEndDate = timestamp
				}
			}
		}
	}

	updateTime, err := unmarshalDatetime(jiraIssue.Fields.Updated)
	if err != nil {
		return Ticket{}, tracerr.Wrap(err)
	}

	createdTime, err := unmarshalDatetime(jiraIssue.Fields.Created)
	if err != nil {
		return Ticket{}, nil
	}

	state := jiraIssue.Fields.Status.Name

	ticket := Ticket{
		Id:         jiraIssue.ID,
		Key:        jiraIssue.Key,
		Title:      jiraIssue.Fields.Summary,
		Type:       jiraIssue.Fields.Type.Name,
		State:      state,
		UpdateTime: updateTime,
		CreateTime: createdTime,

		DevStartDate: devStartDate.Unix(),
		DevEndDate:   devEndDate.Unix(),
	}

	ticket.Transitions = MakeIntervals(ticket, transitions...)
	return ticket, nil
}

func unmarshalDatetime(field jira.Time) (time.Time, error) {
	datetimeRaw, err := field.MarshalJSON()
	if err != nil {
		return BEGINING_OF_TIME, err
	}
	datetime, err := time.Parse(JiraUpdateTimestampFormat, strings.ReplaceAll(string(datetimeRaw), "\"", ""))
	if err != nil {
		return BEGINING_OF_TIME, err
	}
	return datetime, nil
}
