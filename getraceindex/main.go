package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abc123931/keiba-api-aws/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	"gopkg.in/go-playground/validator.v9"
)

var (
	dynamoRegion   string
	dynamoEndpoint string
	validate       *validator.Validate
)

// RaceIndex レース指数用構造体
type RaceIndex struct {
	Name        string `dynamo:"horse_name" json:"horse_name"`
	TotalIndex  string `dynamo:"total_index" json:"total_index"`
	TrainIndex  string `dynamo:"train_index" json:"train_index"`
	StableIndex string `dynamo:"stable_index" json:"stable_index"`
}

// Table dynamo.Table用の構造体
type Table struct {
	Table dynamo.Table
}

// DbConnect Table用のインターフェース
type DbConnect interface {
	get(name string) (raceIndex []RaceIndex, err error)
}

// Request リクエスト用の構造体
type Request struct {
	Name string `json:"name" validate:"required"`
}

// Response レスポンス用の構造体
type Response struct {
	Data  []RaceIndex `json:"data"`
	Error string      `json:"error"`
}

// get DbConnectインターフェースを利用するための関数
func (table *Table) get(name string) (raceIndex []RaceIndex, err error) {
	err = table.Table.Get("race_name", name).All(&raceIndex)
	return
}

// getRaceIndex 検索したいレースの指数のリストを取得する関数
func getRaceIndex(db DbConnect, r events.APIGatewayProxyRequest) *Response {
	request := &Request{}
	response := &Response{Data: []RaceIndex{}}

	err := json.Unmarshal([]byte(r.Body), request)
	if err != nil {
		fmt.Printf("failed unmarshal request: %v", err)
		response.Error = err.Error()
		return response
	}

	validate = validator.New()
	err = validate.Struct(request)

	if err != nil {
		fmt.Printf("validate request: %v", err)
		response.Error = err.Error()
		return response
	}

	response.Data, err = db.get(request.Name)

	if err != nil {
		fmt.Printf("failed get race index: %v", err)
		response.Error = err.Error()
	}

	return response
}

// createResponse レスポンスBodyを生成
func createResponse(response *Response) (responseBody string) {
	json, err := json.Marshal(response)
	if err != nil {
		responseBody = `{"error":{"failed json marshal response"}}`
	}

	responseBody = string(json)
	return
}

// handler ApiGatewayからのリクエストを受けつけ、レスポンスを返却する関数
func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	table := &Table{}
	table.Table = util.ConnectTable("race_index", dynamoRegion, dynamoEndpoint)

	response := createResponse(getRaceIndex(table, r))

	return events.APIGatewayProxyResponse{
		Body:       response,
		StatusCode: 200,
	}, nil
}

func init() {
	util.EnvLoad()
	dynamoRegion = os.Getenv("DYNAMO_REGION")
	dynamoEndpoint = os.Getenv("DYNAMO_ENPOINT")
}

func main() {
	lambda.Start(handler)
}
