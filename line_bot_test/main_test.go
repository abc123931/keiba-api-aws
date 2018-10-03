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

	t.Run("getHorseName異常終了", func(t *testing.T) {
		res := getHorseName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"id": ""}`})
		if res.Error == "" {
			t.Fatal("getHorseName not failed")
		}

		expected := "Key: 'Request.ID' Error:Field validation for 'ID' failed on the 'required' tag"
		if res.Error != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res.Error)
		}
	})

	t.Run("createResponse正常終了", func(t *testing.T) {
		horse := Horse{ID: "test_id", Name: "test_name"}
		response := &Response{Data: horse}
		resBody := createResponse(response)
		expected := `{"data":{"id":"test_id","name":"test_name"},"error":""}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})

	t.Run("createResponseErrorがある場合", func(t *testing.T) {
		response := &Response{Error: "errorだよ"}
		resBody := createResponse(response)
		expected := `{"data":{"id":"","name":""},"error":"errorだよ"}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})
}
