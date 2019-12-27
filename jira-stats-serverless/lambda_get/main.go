package main

import (
	"context"
	"fmt"
	"jira-stats/jira-stats-serverless/jira"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

// Handler is our lambda handler invoked by the `lambda.Start` function call
func mainHandler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	number, err := jira.Fetch()

	result := fmt.Sprintf("Number of fetched Jiras: %d", number)
	if err != nil {
		result = err.Error()
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      201,
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
	lambda.Start(mainHandler)
}
