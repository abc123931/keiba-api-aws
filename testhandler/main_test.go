package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("正常終了", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{}
		res, err := handler(request)
		if err != nil {
			t.Fatal("handler failed")
		}

		if res.StatusCode != 200 {
			t.Fatal("status code is not 200: " + string(res.StatusCode))
		}

		expected := "testだよ"

		assert.Equal(t, expected, res.Body)
	})
}
