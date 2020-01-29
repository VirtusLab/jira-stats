package main

import (
	"context"
	"fmt"
	"github.com/VirtusLab/jira-stats/analyzer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"time"
)

// Handler is our lambda invoked by CloudWatch event
func fetchHandler(ctx context.Context, request events.CloudWatchEvent) (interface{}, error) {
	defer analyzer.TimeTrack(time.Now(), "Jira fetch execution time")

	log.Printf("Jira fetch invoked by: %s at %s\n", request.DetailType, request.Time.Format(time.RFC3339))

	number, err := analyzer.ProcessTickets(100)

	result := fmt.Sprintf("Number of processed Jiras: %d", number)
	if err != nil {
		result = err.Error()
	}

	log.Printf("%s\n", result)
	return result, nil
}

func main() {
	lambda.Start(fetchHandler)
}
