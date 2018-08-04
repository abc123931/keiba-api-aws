package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("getBook正常終了", func(t *testing.T) {
		res, err := getBook(&FakeTable{}, "テストタイトル", "テストカテゴリー")
		if err != nil {
			t.Fatal("getBook failed")
		}

		expected := `{"title":"テストタイトル","category":"テストカテゴリー"}`
		if res != expected {
			t.Fatal("response not same expected: expected is " +
				expected + " response is " + res)
		}
	})

	t.Run("getBook異常終了", func(t *testing.T) {
		_, err := getBook(&FakeTable{}, "", "テストカテゴリー")
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
