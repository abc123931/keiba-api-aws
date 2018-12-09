package main

import (
	"errors"
)

// FakeTable テスト用の構造体
type FakeTable struct{}

func (table *FakeTable) get(name string) (raceIndex []RaceIndex, err error) {
	if name == "" {
		err = errors.New("ValidationException: Comparison type does not exist in DynamoDB")
	}
	raceIndex = append(raceIndex, RaceIndex{Name: "test_horse_name", TotalIndex: "test_total_index"})
	return
}
