package main

import "errors"

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(title string, category string, book *Book) (err error) {
	if title == "" || category == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	book.Title = title
	book.Category = category
	return
}
