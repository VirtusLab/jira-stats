package main

import (
	"context"
	"fmt"
	"github.com/ztrue/tracerr"
	"jira-stats/jira-stats-serverless/analyzer"
	"jira-stats/jira-stats-serverless/analyzer/domain"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func mainHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var result string
	var contentType string

	csv, err := process(request)
	if err != nil {
		tracerr.PrintSourceColor(err)
		result = fmt.Sprintf("Error while generating CSV: %s", err.Error())
		contentType = "application/json"
	}

	result = csv.ToString()
	log.Printf("Generated CSV with: %d rows...", len(csv.Rows)+1)

	contentType = "text/plain"

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            result,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}

	return resp, nil
}

func process(request events.APIGatewayProxyRequest) (*domain.CsvContents, error) {
	params := request.QueryStringParameters

	log.Printf("Path params are: %s", params)

	forceFetch := params["forceFetch"]
	if strings.ToLower(forceFetch) == "true" {
		_, err := analyzer.ProcessTickets(100)
		return &domain.CsvContents{}, err
	}

	startDateString := params["startDate"]
	startDate, err := time.Parse(domain.DayFormat, startDateString)
	if err != nil {
		return &domain.CsvContents{}, tracerr.Wrap(err)
	}

	endDateString := params["endDate"]
	endDate, err := time.Parse(domain.DayFormat, endDateString)
	if err != nil {
		return &domain.CsvContents{}, tracerr.Wrap(err)
	}

	csv, err := analyzer.GetCsv(startDate, endDate)
	if err != nil {
		return &domain.CsvContents{}, tracerr.Wrap(err)
	}

	return csv, nil
}

func main() {
	lambda.Start(mainHandler)
}
