package analyzer

import (
	"fmt"
	"github.com/VirtusLab/jira-stats/analyzer/domain"
	"github.com/andygrunwald/go-jira"
	"github.com/ztrue/tracerr"
	"log"
	"time"
)

const GoogleSpreadsheetFormat = "2006-01-02 15:04:05"

const MaxBatchSize = 400

// Fetches data from Jira and stores it in DB
func ProcessTickets(batchCount int) (int, error) {
	if batchCount > MaxBatchSize {
		return -1, fmt.Errorf("requested batch size [%d] bigger than allowed limit [%d]", batchCount, MaxBatchSize)
	}

	count, err := processLoop(batchCount)
	if err != nil {
		return -1, tracerr.Wrap(err)
	}

	if count < batchCount {
		log.Printf("Read all the issues up to date...")
	} else {
		log.Printf("More issues most likely to be read on next execution...")
	}

	return count, nil
}

func processLoop(batchCount int) (int, error) {
	// gets last update to figure out where to start with fetching
	lastUpdate, err := getLastUpdate()
	if err != nil {
		return -1, tracerr.Wrap(err)
	}
	log.Printf("Last update is: %s\n", lastUpdate.Format(time.RFC3339))

	// fetches tickets
	jiraTickets, err := fetch(lastUpdate, batchCount)
	if err != nil {
		return -1, tracerr.Wrap(err)
	}

	// converts Jira issues to model
	tickets, err := transformToModel(jiraTickets)
	if err != nil {
		return -1, err
	}

	// stores in db
	mostRecentUpdate, err := storeTickets(tickets)
	if err != nil {
		return -1, err
	}

	// updates update time for next round
	processTicketsNo := len(tickets)

	err = storeLastUpdate(mostRecentUpdate)
	if err != nil {
		return processTicketsNo, tracerr.Wrap(err)
	}

	log.Printf("Processed %d tickets...\n", processTicketsNo)
	return processTicketsNo, nil
}

func transformToModel(jiraTickets []jira.Issue) (tickets []domain.Ticket, err error) {
	defer timeTrackParams(time.Now(), "Converting tickets to model", map[string]string{"number": string(len(tickets))})

	tickets, err = BuildModel(jiraTickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func storeTickets(tickets []domain.Ticket) (lastUpdateTime time.Time, err error) {
	defer timeTrackParams(time.Now(), "Storing tickets", map[string]string{"number": string(len(tickets))})

	mostRecentUpdate := domain.BeginingOfTime
	// stores new model
	for _, ticket := range tickets {
		err = store(ticket)
		if err != nil {
			return time.Time{}, err
		}

		updateTime := ticket.UpdateTime

		if mostRecentUpdate.Before(updateTime) {
			mostRecentUpdate = updateTime
		}
	}

	return mostRecentUpdate, nil
}

// transforms analyzer tickets to model
func BuildModel(jiraIssues []jira.Issue) ([]domain.Ticket, error) {
	domainTickets := make([]domain.Ticket, 0)
	for _, issue := range jiraIssues {
		domainTicket, err := domain.JiraToDomain(issue)
		if err != nil {
			return nil, err
		}

		updateString := domainTicket.UpdateTime.Format(time.RFC3339)

		log.Printf("Ticket: %s, key: %s, no of changes: %d (updated at %s)",
			domainTicket.Id, domainTicket.Key, len(domainTicket.ChangelogEntries), updateString)

		domainTickets = append(domainTickets, domainTicket)
	}

	return domainTickets, nil
}
