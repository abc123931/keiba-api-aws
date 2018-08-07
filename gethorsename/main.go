package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
)

var (
	dynamoRegion   string
	dynamoEndpoint string
)

// Horse 馬名検索用の構造体
type Horse struct {
	Category string `dynamo:"category" json:"category"`
	HorseID  string `dynamo:"horse_id" json:"id"`
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
	Category  string `json:category`
	HorseName string `json:horse_name`
}

// Env_load 開発環境ようにdotenvを読み込む関数
func Env_load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("start lambda function at prod")
	}
}

// get DbConnectインターフェースを利用するための関数
func (table *Table) get(category string, name string) (horses []Horse, err error) {
	horses = []Horse{}
	err = table.Table.Get("category", category).
		Filter("contains($, ?)", "name", name).
		All(&horses)

	return
}

// getHorseName 検索したい馬のリストを取得する関数
func getHorseName(db DbConnect, category string, horseName string) (response string, err error) {
	horses, err := db.get(category, horseName)
	if err != nil {
		fmt.Printf("failed get horse names: %v", err)
		return
	}

	json, err := json.Marshal(horses)
	if err != nil {
		log.Println(err)
		return
	}
	response = string(json)
	return

}

// handler ApiGatewayからのリクエストを受けつけ、レスポンスを返却する関数
func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String(dynamoRegion),
		Endpoint: aws.String(dynamoEndpoint),
	})

	table := &Table{db.Table("Horses")}

	request := &Request{}
	err := json.Unmarshal([]byte(r.Body), request)

	response := ""

	if err != nil {
		response = fmt.Sprintf("cannot encode request json: %s", err)
	} else {
		response, err = getHorseName(table, request.Category, request.HorseName)
		if err != nil {
			response = fmt.Sprintf("cannot get horse name: %s", err)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       response,
		StatusCode: 200,
	}, nil
}

func init() {
	Env_load()
	dynamoRegion = os.Getenv("DYNAMO_REGION")
	dynamoEndpoint = os.Getenv("DYNAMO_ENPOINT")
}

func main() {
	lambda.Start(handler)
}
