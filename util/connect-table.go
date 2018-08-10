package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func ConnectTable(tableName string, dynamoRegion string, dynamoEndpoint string) dynamo.Table {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String(dynamoRegion),
		Endpoint: aws.String(dynamoEndpoint),
	})

	table := db.Table(tableName)

	return table
}
