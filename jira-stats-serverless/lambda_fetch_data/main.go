package main

import (
	"context"
	"fmt"
	"jira-stats/jira-stats-serverless/jira"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func fetchHandler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	number, err := jira.Fetch()

	result := fmt.Sprintf("Number of fetched Jiras: %d", number)
	if err != nil {
		result = err.Error()
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            result,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "lambda_get-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(fetchHandler)
}
