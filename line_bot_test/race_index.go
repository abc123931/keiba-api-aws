package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

// RaceIndexRequest レースインデックスリクエスト構造体
type RaceIndexRequest struct {
	Name string `json:"name"`
}

// RaceIndexData raceIndexのレスポンスのdataの構造体
type RaceIndexData struct {
	Data []RaceIndex `json:"data"`
}

// RaceIndexs ソートするための構造体のスライス
type RaceIndexs []RaceIndex

// 以下Sortインターフェースを継承する為に必要
func (r RaceIndexs) Len() int {
	return len(r)
}

func (r RaceIndexs) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// ここでなにでソートするか決める
func (r RaceIndexs) Less(i, j int) bool {
	return r[i].TotalIndex > r[j].TotalIndex
}

// RaceIndex レースインデックス用構造体
type RaceIndex struct {
	Name        string `dynamo:"horse_name" json:"horse_name"`
	TotalIndex  string `dynamo:"total_index" json:"total_index"`
	TrainIndex  string `dynamo:"train_index" json:"train_index"`
	StableIndex string `dynamo:"stable_index" json:"stable_index"`
}

// requestRaceIndexAPI RaceIndexApiクライアント
func requestRaceIndexAPI(name string) []RaceIndex {
	values, err := json.Marshal(RaceIndexRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/raceindex", "application/json", bytes.NewBuffer(values))

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	raceIndexData := &RaceIndexData{}

	err = json.Unmarshal(body, raceIndexData)

	if err != nil {
		fmt.Printf("failed unmarshal request: %v", err)
		log.Fatal(err)
	}

	var raceIndexs RaceIndexs = raceIndexData.Data

	sort.Sort(raceIndexs)

	return raceIndexs
}
