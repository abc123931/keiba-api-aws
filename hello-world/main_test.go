package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	request := events.APIGatewayProxyRequest{}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello World",
	}

	response, err := Handler(request)

	assert.Equal(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)

}
