package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("getHorseName正常終了", func(t *testing.T) {
		res := getHorseName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"id": "test_id"}`})
		if res.Error != "" {
			t.Fatal("getHorseName failed: ", res.Error)
		}

		horse := Horse{ID: "test_id", Name: "test_name"}
		expected := &Response{Data: horse, Error: ""}
		if !reflect.DeepEqual(*res, *expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", *expected, *res)
		}
	})
}
