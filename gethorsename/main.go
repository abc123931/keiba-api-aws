package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abc123931/keiba-api-aws/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
)

var (
	dynamoRegion   string
	dynamoEndpoint string
)

// Horse 馬名検索用の構造体
type Horse struct {
	Category string `dynamo:"category" json:"category"`
	HorseID  string `dynamo:"id" json:"id"`
	Name     string `dynamo:"name" json:"name"`
}

// Table dynamo.Table用の構造体
type Table struct {
	Table dynamo.Table
}

// DbConnect Table用のインターフェース
type DbConnect interface {
	get(category string, name string) (horses []Horse, err error)
}

// Request リクエスト用の構造体
type Request struct {
	Category  string `json:"category"`
	HorseName string `json:"horse_name"`
}

// Response レスポンス用の構造体
type Response struct {
	Data  []Horse `json:"data"`
	Error string  `json:"error"`
}

// get DbConnectインターフェースを利用するための関数
func (table *Table) get(category string, name string) (horses []Horse, err error) {
	horses = []Horse{}
	err = table.Table.Get("category", category).
		Range("name", dynamo.BeginsWith, name).
		All(&horses)

	return
}

// getHorseName 検索したい馬のリストを取得する関数
func getHorseName(db DbConnect, r events.APIGatewayProxyRequest) *Response {
	request := &Request{}
	response := &Response{}

	err := json.Unmarshal([]byte(r.Body), request)
	if err != nil {
		fmt.Printf("failed unmarshal request: %v", err)
		response.Error = err.Error()
		return response
	}

	response.Data, err = db.get(request.Category, request.HorseName)

	if err != nil {
		fmt.Printf("failed get horse names: %v", err)
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
	table.Table = util.ConnectTable("search_horses", dynamoRegion, dynamoEndpoint)

	response := createResponse(getHorseName(table, r))

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