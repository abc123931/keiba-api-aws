package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(category string, name string) (races []Race, err error) {
	races = []Race{}
	if category == "" || name == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	races = append(races, Race{category, "test_race_id", name})
	return
}
