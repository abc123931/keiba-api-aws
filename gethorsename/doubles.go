package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(category string, name string) (horses []Horse, err error) {
	horses = []Horse{}
	if category == "" || name == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	horses = append(horses, Horse{category, "test_horse_id", name})
	return
}
