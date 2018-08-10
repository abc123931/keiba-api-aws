package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("getHorseName正常終了", func(t *testing.T) {
		res := getHorseName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"category": "horse", "horse_name": "test_horse_name"}`})
		if res.Error != "" {
			t.Fatal("getHorseName failed: ", res.Error)
		}

		horses := []Horse{}
		horses = append(horses, Horse{Category: "horse", HorseID: "test_horse_id", Name: "test_horse_name"})
		expected := &Response{Data: horses, Error: ""}
		if !reflect.DeepEqual(*res, *expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", *expected, *res)
		}
	})

	t.Run("getHorseName異常終了", func(t *testing.T) {
		res := getHorseName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"category": "", "horse_name": "test_horse_name"}`})
		if res.Error == "" {
			t.Fatal("getHorseName not failed")
		}

		expected := "ValidationException: Comparison type does not exist in DynamoDB"
		if res.Error != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res.Error)
		}
	})

	t.Run("createResponse正常終了", func(t *testing.T) {
		horses := []Horse{}
		horses = append(horses, Horse{Category: "horse", HorseID: "test_horse_id", Name: "test_horse_name"})
		response := &Response{Data: horses}
		resBody := createResponse(response)
		expected := `{"data":[{"category":"horse","id":"test_horse_id","name":"test_horse_name"}],"error":""}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})

	t.Run("createResponseErrorがある場合", func(t *testing.T) {
		response := &Response{Error: "errorだよ"}
		resBody := createResponse(response)
		expected := `{"data":null,"error":"errorだよ"}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})
}
