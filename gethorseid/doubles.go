package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(id string) (horse Horse, err error) {
	horse = Horse{}
	if id == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	horse = Horse{ID: "test_id", Name: "test_name"}
	return
}
