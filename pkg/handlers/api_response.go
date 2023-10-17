package handlers

import (
	"encoding/json"
	"log"
	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	res := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type":"application/json"}}
	res.StatusCode = status

	marshalledBody, err := json.Marshal(body)

	if err != nil {
		log.Fatal(err)
	}

	res.Body = string(marshalledBody)
	return &res, nil
}