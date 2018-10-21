package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(id string) (courseResult CourseResult, err error) {
	courseResult = CourseResult{}
	if id == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	courseResult = CourseResult{ID: "test_id"}
	return
}
