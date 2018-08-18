package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("getRaceName正常終了", func(t *testing.T) {
		res := getRaceName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"category": "race", "race_name": "test_race_name"}`})
		if res.Error != "" {
			t.Fatal("getRaceName failed: ", res.Error)
		}

		races := []Race{}
		races = append(races, Race{Category: "race", RaceID: "test_race_id", Name: "test_race_name"})
		expected := &Response{Data: races, Error: ""}
		if !reflect.DeepEqual(*res, *expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", *expected, *res)
		}
	})

	t.Run("getRaceName異常終了", func(t *testing.T) {
		res := getRaceName(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"category": "", "race_name": "test_race_name"}`})
		if res.Error == "" {
			t.Fatal("getRaceName not failed")
		}

		expected := "Key: 'Request.Category' Error:Field validation for 'Category' failed on the 'required' tag"
		if res.Error != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res.Error)
		}
	})

	t.Run("createResponse正常終了", func(t *testing.T) {
		races := []Race{}
		races = append(races, Race{Category: "race", RaceID: "test_race_id", Name: "test_race_name"})
		response := &Response{Data: races}
		resBody := createResponse(response)
		expected := `{"data":[{"category":"race","id":"test_race_id","name":"test_race_name"}],"error":""}`

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
