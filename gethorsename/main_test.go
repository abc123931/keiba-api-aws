package main

import (
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("getHorseName正常終了", func(t *testing.T) {
		res, err := getHorseName(&FakeTable{}, "horse", "test_horse_name")
		if err != nil {
			t.Fatal("getHorseName failed", err)
		}

		expected := []string{`{"category":"horse","id":"test_horse_id","name":"test_horse_name"}`}
		if reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})

	t.Run("getHorseName異常終了", func(t *testing.T) {
		_, err := getHorseName(&FakeTable{}, "", "test_horse_name")
		if err == nil {
			t.Fatal("getBook not failed")
		}

		expected := "ValidationException: Comparison type does not exist in DynamoDB"
		if err.Error() != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + err.Error())
		}
	})

}
