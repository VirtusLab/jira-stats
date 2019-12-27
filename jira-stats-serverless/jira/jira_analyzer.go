package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"time"
)

func Fetch() (int, error) {

	tp := jira.BasicAuthTransport{
		Username: "xxxxx",
		Password: "xxxxx",
	}

	lastUpdate, err := getLastUpdate()
	if err != nil {
		return -1, tracerr.Wrap(err)
	}
	fmt.Printf("Last update is: %s", lastUpdate.Format(time.RFC3339))

	client, _ := jira.NewClient(tp.Client(), "https://jira.adstream.com")

	jqlQuery :=
		"(project in (\"Traffic & Ordering\", \"Amazing Delivery\", ROB) OR\n" +
			"\t kanban = \"Traffic & Ordering\" OR\n" +
			"\tlabels in (traffic-external, traffic-team)\n" +
			") AND \n" +
			"project != NIR AND \n" +
			"NOT (project = DL AND status = Closed) AND \n" +
			"updated > 2019-10-01 " +
			"ORDER BY Rank ASC"

	fmt.Println(jqlQuery)

	searchOpts := jira.SearchOptions{MaxResults: 5, Expand: "changelog"}
	issues, _, _ := client.Issue.Search(jqlQuery, &searchOpts)

	fmt.Println(issues)

	for _, issue := range issues {

		domainTicket := jiraToDomain(issue)
		err := store(domainTicket)
		if err != nil {
			return -1, err
		}

		for _, historyItem := range issue.Changelog.Histories {
			for _, changeItem := range historyItem.Items {
				if changeItem.Field == "status" {
					fmt.Printf("%s changed status %s -> %s\n", issue.Key, changeItem.FromString, changeItem.ToString)
				}
			}
		}
	}

	issue, _, err := client.Issue.Get("PTR-19", nil)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)

	err = storeLastUpdate(time.Now())
	if err != nil {
		return -1, tracerr.Wrap(err)
	}

	return len(issues), nil
}
