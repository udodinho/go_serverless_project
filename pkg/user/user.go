package user

import (
	"encoding/json"
	"errors"
	"github/udodinho/go-serverless/pkg/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

var (
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnmashalRecord = "failed to unmarshal record"
	ErrorInvalidUserData = "invalid user data"
	ErrorInvalidEmail = "invalid email"
	ErrorCouldNotMarshalItem = "could not marshal item"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorUserAlreadyExist = "user already exist"
	ErrorUserDoesNotExist = "user does not exist"
)

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	res, err := dynaClient.GetItem(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)

	err = dynamodbattribute.UnmarshalMap(res.Item, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmashalRecord)
	}

	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	res, err := dynaClient.Scan(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)

	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmashalRecord)
	}

	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var user User

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	exist, _ := FetchUser(user.Email, tableName, dynaClient)
	
	if exist != nil && len(exist.Email)!= 0 {
		return nil, errors.New(ErrorFailedToUnmashalRecord)
	}

	newUser, err := dynamodbattribute.MarshalMap(user)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: newUser,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var user User

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	exist, _ := FetchUser(user.Email, tableName, dynaClient)

	if exist != nil && len(exist.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	newUser, err := dynamodbattribute.MarshalMap(user)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: newUser,
		TableName: &tableName,
	}
	_, err = dynaClient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil

}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)

	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}