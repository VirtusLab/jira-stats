package main

import (
	"github.com/ztrue/tracerr"
	"jira-stats/jira-stats-serverless/analyzer"
	"jira-stats/jira-stats-serverless/analyzer/domain"
	"log"
	"time"
)

func main() {
	//_, err := analyzer.ProcessTickets(100)
	//if err != nil {
	//	tracerr.PrintSourceColor(err)
	//}
	//
	start, _ := time.Parse(domain.DayFormat, "2019-10-01")
	end, _ := time.Parse(domain.DayFormat, "2019-12-31")
	csv, err := analyzer.GetCsv(start, end)
	if err != nil {
		tracerr.PrintSourceColor(err)
	}

	log.Printf("%s", csv.ToString())
}
