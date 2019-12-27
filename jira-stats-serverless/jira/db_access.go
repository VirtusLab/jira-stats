package jira

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ztrue/tracerr"
	"log"
	"time"
	//"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	//"os"
)

var EMPTY_DATE = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func store(ticket Ticket) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	item, err := dynamodbattribute.MarshalMap(ticket)
	if err != nil {
		return tracerr.Wrap(err)
	}

	input := dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Ticket"),
	}
	_, err = svc.PutItem(&input)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func storeLastUpdate(updateTime time.Time) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	prevUpdate, err := getLastUpdate()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if prevUpdate != EMPTY_DATE {
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":w": {
					S: aws.String(updateTime.Format(time.RFC3339)),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"ConfigName": {
					S: aws.String("LastUpdate"),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			TableName:        aws.String("Config"),
			UpdateExpression: aws.String("set ConfigValue = :w"),
		}
		output, err := svc.UpdateItem(input)
		if err != nil {
			return tracerr.Wrap(err)
		}
		log.Printf("Output is: %s", output)
	} else {
		configItem := ConfigItem{
			ConfigName:  "LastUpdate",
			ConfigValue: updateTime.Format(time.RFC3339),
		}

		item, err := dynamodbattribute.MarshalMap(configItem)
		if err != nil {
			return tracerr.Wrap(err)
		}

		input := dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String("Config"),
		}

		_, err = svc.PutItem(&input)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return err
}

func getLastUpdate() (time.Time, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConfigName": {
				S: aws.String("LastUpdate"),
			},
		},

		TableName: aws.String("Config"),
	})
	if err != nil {
		return time.Now(), tracerr.Wrap(err)
	}

	configResult := ConfigItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &configResult)
	if err != nil {
		return time.Now(), tracerr.Wrap(err)
	}

	lastUpdate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	if configResult.ConfigValue != "" {
		lastUpdate, err = time.Parse(time.RFC3339, configResult.ConfigValue)
		if err != nil {
			return time.Now(), tracerr.Wrap(err)
		}
	}

	return lastUpdate, nil
}
