package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("getRaceIndex正常終了", func(t *testing.T) {
		res := getRaceIndex(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"name": "test_name"}`})
		if res.Error != "" {
			t.Fatal("getRaceIndex failed: ", res.Error)
		}

		raceIndex := []RaceIndex{}
		raceIndex = append(raceIndex, RaceIndex{Name: "test_horse_name", TotalIndex: "test_total_index"})
		expected := &Response{Data: raceIndex, Error: ""}
		if !reflect.DeepEqual(*res, *expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", *expected, *res)
		}
	})

	t.Run("getRaceIndex異常終了", func(t *testing.T) {
		res := getRaceIndex(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"name": ""}`})
		if res.Error == "" {
			t.Fatal("getRaceIndex not failed")
		}

		expected := "Key: 'Request.Name' Error:Field validation for 'Name' failed on the 'required' tag"
		if res.Error != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res.Error)
		}
	})

	t.Run("createResponse正常終了", func(t *testing.T) {
		raceIndex := []RaceIndex{}
		raceIndex = append(raceIndex, RaceIndex{Name: "test_name", TotalIndex: "83.4", TrainIndex: "-", StableIndex: "6.7"})
		response := &Response{Data: raceIndex}
		resBody := createResponse(response)
		expected := `{"data":[{"horse_name":"test_name","total_index":"83.4","train_index":"-","stable_index":"6.7"}],"error":""}`

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
