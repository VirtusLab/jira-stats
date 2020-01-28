package domain

import (
	//"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"strings"
	"time"
)

var BeginingOfTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
var EndOfTime = time.Date(9999, 12, 31, 23, 59, 59, 999, time.UTC)

const JiraTimestampFormat = "2006-01-02T15:04:05.000-0700"
const JiraUpdateTimestampFormat = "2006-01-02T15:04:05-0700"

const JiraFilterFormat = "2006-01-02 15:04"
const DayFormat = "2006-01-02"

type Ticket struct {
	Id               string
	Key              string
	State            string
	Type             string
	Title            string
	ChangelogEntries []ChangelogEntry
	UpdateTime       time.Time
	CreateTime       time.Time

	DevStartDate int64
	DevEndDate   int64
}

func (t *Ticket) Project() string {
	dashIdx := strings.LastIndex(t.Key, "-")
	return t.Key[0:dashIdx]
}

type ChangelogEntry struct {
	Id      string
	Author  string
	Created time.Time
	Changes []Change
}

func CreateChangelogEntry(id string, author string, timestamp time.Time, changes ...Change) ChangelogEntry {
	return ChangelogEntry{
		Id:      id,
		Author:  author,
		Created: timestamp,
		Changes: changes,
	}
}

type Change struct {
	Field string
	From  string
	To    string
}

func CreateChange(field string, from string, to string) Change {
	return Change{
		Field: field,
		From:  from,
		To:    to,
	}
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

	changeslogEntries := make([]ChangelogEntry, 0)

	devStartDate := EndOfTime
	devEndDate := BeginingOfTime

	for _, historyItem := range jiraIssue.Changelog.Histories {
		changes := make([]Change, 0)

		for _, changeItem := range historyItem.Items {
			if strings.ToLower(changeItem.Field) == "status" {

				timestamp, err := time.Parse(JiraTimestampFormat, historyItem.Created)
				if err != nil {
					return Ticket{}, tracerr.Wrap(err)
				}

				changes = append(changes, Change{
					From:  changeItem.FromString,
					To:    changeItem.ToString,
					Field: changeItem.Field,
				})

				if changeItem.ToString == "In Development" && devStartDate.After(timestamp) {
					devStartDate = timestamp
				}

				if changeItem.FromString == "In Development" && devEndDate.Before(timestamp) {
					devEndDate = timestamp
				}
			}

			changeLogEntry, err := buildChangelogEntry(historyItem, changes)
			changeslogEntries = append(changeslogEntries, changeLogEntry)
			if err != nil {
				return Ticket{}, tracerr.Wrap(err)
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

		ChangelogEntries: changeslogEntries,

		DevStartDate: devStartDate.Unix(),
		DevEndDate:   devEndDate.Unix(),
	}

	return ticket, nil
}

func buildChangelogEntry(historyItem jira.ChangelogHistory, changes []Change) (ChangelogEntry, error) {
	changeTime, err := historyItem.CreatedTime()
	if err != nil {
		return ChangelogEntry{}, tracerr.Wrap(err)
	}

	return ChangelogEntry{
		Id:      historyItem.Id,
		Author:  historyItem.Author.Name,
		Created: changeTime,
		Changes: changes,
	}, nil
}

func unmarshalDatetime(field jira.Time) (time.Time, error) {
	datetimeRaw, err := field.MarshalJSON()
	if err != nil {
		return BeginingOfTime, err
	}
	datetime, err := time.Parse(JiraUpdateTimestampFormat, strings.ReplaceAll(string(datetimeRaw), "\"", ""))
	if err != nil {
		return BeginingOfTime, err
	}
	return datetime, nil
}
