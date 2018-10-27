package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("getCourseResult正常終了", func(t *testing.T) {
		res := getCourseResult(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"id": "test_id"}`})
		if res.Error != "" {
			t.Fatal("getCourseResult failed: ", res.Error)
		}

		courseResult := CourseResult{ID: "test_id"}
		expected := &Response{Data: courseResult, Error: ""}
		if !reflect.DeepEqual(*res, *expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", *expected, *res)
		}
	})

	t.Run("getCourseResult異常終了", func(t *testing.T) {
		res := getCourseResult(&FakeTable{}, events.APIGatewayProxyRequest{Body: `{"id": ""}`})
		if res.Error == "" {
			t.Fatal("getCourseResult not failed")
		}

		expected := "Key: 'Request.ID' Error:Field validation for 'ID' failed on the 'required' tag"
		if res.Error != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res.Error)
		}
	})

	t.Run("createResponse正常終了", func(t *testing.T) {
		courseResult := CourseResult{ID: "test_id", SapporoTurf: "(0-1-0-0)"}
		response := &Response{Data: courseResult}
		resBody := createResponse(response)
		expected := `{"data":{"id":"test_id","sapporo_turf":"(0-1-0-0)","hakodate_turf":"","fukushima_turf":"","nigata_turf":"","tokyo_turf":"","nakayama_turf":"","tyukyo_turf":"","kyoto_turf":"","hanshin_turf":"","kokura_turf":"","1000_turf":"","1200_turf":"","1400_turf":"","1600_turf":"","1800_turf":"","2000_turf":"","2200_turf":"","2400_turf":"","2500_turf":"","3000_turf":"","3200_turf":"","3600_turf":""},"error":""}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})

	t.Run("createResponseErrorがある場合", func(t *testing.T) {
		response := &Response{Error: "errorだよ"}
		resBody := createResponse(response)
		expected := `{"data":{"id":"","sapporo_turf":"","hakodate_turf":"","fukushima_turf":"","nigata_turf":"","tokyo_turf":"","nakayama_turf":"","tyukyo_turf":"","kyoto_turf":"","hanshin_turf":"","kokura_turf":"","1000_turf":"","1200_turf":"","1400_turf":"","1600_turf":"","1800_turf":"","2000_turf":"","2200_turf":"","2400_turf":"","2500_turf":"","3000_turf":"","3200_turf":"","3600_turf":""},"error":"errorだよ"}`

		if resBody != expected {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, resBody)
		}
	})
}
