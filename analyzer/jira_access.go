package analyzer

import (
	"encoding/json"
	"fmt"
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"log"
	"os"
	"time"
)

// fetches tickets from analyzer
func fetch(updatedSince time.Time, batchCount int) ([]jira.Issue, error) {
	defer TimeTrack(time.Now(), "Fetching Jira issues")

	tp, err := jiraAuth()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	client, err := jira.NewClient(tp.Client(), "https://jira.adstream.com")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	jqlQuery :=
		fmt.Sprintf("(project in (\"Traffic & Ordering\", \"Amazing Delivery\", ROB) OR "+
			"kanban = \"Traffic & Ordering\" OR "+
			"labels in (traffic-external, traffic-team) "+
			") AND "+
			"project != NIR AND "+
			"NOT (project = DL AND status = Closed) AND "+
			"updated >= \"%s\" "+
			"ORDER BY updated ASC",

			updatedSince.Format(domain.JiraFilterFormat),
		)

	log.Printf("Jira query used: %s\n", jqlQuery)

	searchOpts := jira.SearchOptions{MaxResults: batchCount, Expand: "changelog"}
	issues, _, err := client.Issue.Search(jqlQuery, &searchOpts)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return issues, nil
}

// Builds analyzer client (fetching creds either from env vars or AWS Secret Manager)
func jiraAuth() (jira.BasicAuthTransport, error) {

	if os.Getenv("JIRA_USER") != "" {
		log.Printf("Fetching creds from local vars...")
		return jira.BasicAuthTransport{
			Username: os.Getenv("JIRA_USER"),
			Password: os.Getenv("JIRA_PASSWORD"),
		}, nil
	} else {
		log.Printf("Fetching creds from Secret Manager...")
		secrets, err := RetrieveSecrets()
		if err != nil {
			return jira.BasicAuthTransport{}, tracerr.Wrap(err)
		}

		type Creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		creds := Creds{}
		err = json.Unmarshal(secrets, &creds)
		if err != nil {
			return jira.BasicAuthTransport{}, tracerr.Wrap(err)
		}

		return jira.BasicAuthTransport{
			Username: creds.Username,
			Password: creds.Password,
		}, nil
	}
}
