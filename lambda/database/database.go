package database

import (
	"lambda_func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	TABLE_NAME = "userTable"
)

type DynamoDBClient struct {
	dbStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		dbStore: db,
	}
}

func (d DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := d.dbStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (d DynamoDBClient) InsertUser(user types.RegisterUser) error {
	// assemble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.Password),
			},
		},
	}
	// insert item
	_, err := d.dbStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}
