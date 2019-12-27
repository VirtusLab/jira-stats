package jira

import (
	"github.com/andygrunwald/go-jira"
)

type Ticket struct {
	Id    string
	Title string
}

type ConfigItem struct {
	ConfigName  string
	ConfigValue string
}

func jiraToDomain(jiraIssue jira.Issue) Ticket {
	return Ticket{Id: jiraIssue.Key, Title: jiraIssue.Fields.Summary}
}
