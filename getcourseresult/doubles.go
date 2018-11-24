package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(name string) (courseResult CourseResult, err error) {
	courseResult = CourseResult{}
	if name == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	courseResult = CourseResult{Name: "test_name"}
	return
}
