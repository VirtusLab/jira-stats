package integration

import (
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"log"
	"os"
	"testing"
)

func TestJiraConnection(t *testing.T) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}

	client, _ := jira.NewClient(tp.Client(), "https://jira.adstream.com")

	issue, _, err := client.Issue.Get("PTR-19", nil)
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Error found: %s", err.Error())
	}

	if issue.Key != "PTR-19" {
		t.Errorf("Incorrect issue key: %s", issue.Key)
	}

	if issue.Self != "https://jira.adstream.com/rest/api/2/issue/137941" {
		t.Errorf("Incorrect self link: %s", issue.Self)
	}

	log.Printf("Issue: %+v", issue)
}
