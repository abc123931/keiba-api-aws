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

// CourseResult コース成績用構造体
type CourseResult struct {
	Name                  string `dynamo:"name" json:"name"`
	SapporoTurf           string `dynamo:"sapporo_turf" json:"sapporo_turf"`
	HakodateTurf          string `dynamo:"hakodate_turf" json:"hakodate_turf"`
	HukushimaTurf         string `dynamo:"hukushima_turf" json:"hukushima_turf"`
	NigataTurf            string `dynamo:"nigata_turf" json:"nigata_turf"`
	TokyoTurf             string `dynamo:"tokyo_turf" json:"tokyo_turf"`
	NakayamaTurf          string `dynamo:"nakayama_turf" json:"nakayama_turf"`
	TyukyoTurf            string `dynamo:"tyukyo_turf" json:"tyukyo_turf"`
	KyotoTurf             string `dynamo:"kyoto_turf" json:"kyoto_turf"`
	HanshinTurf           string `dynamo:"hanshin_turf" json:"hanshin_turf"`
	KokuraTurf            string `dynamo:"kokura_turf" json:"kokura_turf"`
	ThousandTurf          string `dynamo:"1000_turf" json:"1000_turf"`
	TwelveHundredTurf     string `dynamo:"1200_turf" json:"1200_turf"`
	FourteenHundredTurf   string `dynamo:"1400_turf" json:"1400_turf"`
	SixteenHundredTurf    string `dynamo:"1600_turf" json:"1600_turf"`
	EighteenHundredTurf   string `dynamo:"1800_turf" json:"1800_turf"`
	TwoThousandTurf       string `dynamo:"2000_turf" json:"2000_turf"`
	TwentyTwoHundredTurf  string `dynamo:"2200_turf" json:"2200_turf"`
	TwentyFourHundredTurf string `dynamo:"2400_turf" json:"2400_turf"`
	TwentyFiveHundredTurf string `dynamo:"2500_turf" json:"2500_turf"`
	ThreeThousandTurf     string `dynamo:"3000_turf" json:"3000_turf"`
	ThirtyTwoHundredTurf  string `dynamo:"3200_turf" json:"3200_turf"`
	ThirtySixHundredTurf  string `dynamo:"3600_turf" json:"3600_turf"`
	RyoTurf               string `dynamo:"ryo_turf" json:"ryo_turf"`
	YayaomoTurf           string `dynamo:"yayaomo_turf" json:"yayaomo_turf"`
	OmoTurf               string `dynamo:"omo_turf" json:"omo_turf"`
	HuryoTurf             string `dynamo:"huryo_turf" json:"huryo_turf"`
}

// Table dynamo.Table用の構造体
type Table struct {
	Table dynamo.Table
}

// DbConnect Table用のインターフェース
type DbConnect interface {
	get(name string) (courseResult CourseResult, err error)
}

// Request リクエスト用の構造体
type Request struct {
	Name string `json:"name" validate:"required"`
}

// Response レスポンス用の構造体
type Response struct {
	Data  CourseResult `json:"data"`
	Error string       `json:"error"`
}

// get DbConnectインターフェースを利用するための関数
func (table *Table) get(name string) (courseResult CourseResult, err error) {
	courseResult = CourseResult{}
	err = table.Table.Get("horse_name", name).One(&courseResult)

	return
}

// getCourseResult 検索したい馬のリストを取得する関数
func getCourseResult(db DbConnect, r events.APIGatewayProxyRequest) *Response {
	request := &Request{}
	response := &Response{Data: CourseResult{}}

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
	table.Table = util.ConnectTable("course_results", dynamoRegion, dynamoEndpoint)

	response := createResponse(getCourseResult(table, r))

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
