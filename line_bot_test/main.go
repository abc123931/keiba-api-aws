package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/abc123931/keiba-api-aws/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	channelSecret string
	channelToken  string
)

// LineMessage lineから送信されたMessageの構造体
type LineMessage struct {
	Message    *linebot.TextMessage
	ReplyToken string
	Status     int
}

// HorseNameRequest 馬名のリクエスト構造体
type HorseNameRequest struct {
	Category  string `json:"category"`
	HorseName string `json:"horse_name"`
}

// CourseResultRequest 馬名のリクエスト構造体
type CourseResultRequest struct {
	Name string `json:"name"`
}

// HorseNameData gethorsenameのレスポンスのdataの構造体
type HorseNameData struct {
	Data []HorseNameResponse `json:"data"`
}

// CourseResultData courseresultのレスポンスのdataの構造体
type CourseResultData struct {
	Data CourseResult `json:"data"`
}

// HorseNameResponse 馬名のリクエスト構造体
type HorseNameResponse struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

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

var courseResultName = map[string]string{
	"SapporoTurf":   "札幌成績",
	"HakodateTurf":  "函館成績",
	"HukushimaTurf": "福島成績",
	"NigataTurf":    "新潟成績",
	"TokyoTurf":     "東京成績",
	"NakayamaTurf":  "中山成績",
	"TyukyoTurf":    "中京成績",
	"KyotoTurf":     "京都成績",
	"HanshinTurf":   "阪神成績",
	"KokuraTurf":    "小倉成績",
}

var distanceResultName = map[string]string{
	"ThousandTurf":          "芝1000m",
	"TwelveHundredTurf":     "芝1200m",
	"FourteenHundredTurf":   "芝1400m",
	"SixteenHundredTurf":    "芝1600m",
	"EighteenHundredTurf":   "芝1800m",
	"TwoThousandTurf":       "芝2000m",
	"TwentyTwoHundredTurf":  "芝2200m",
	"TwentyFourHundredTurf": "芝2400m",
	"TwentyFiveHundredTurf": "芝2500m",
	"ThreeThousandTurf":     "芝3000m",
	"ThirtyTwoHundredTurf":  "芝3200m",
	"ThirtySixHundredTurf":  "芝3600m",
}

var babaResultName = map[string]string{
	"RyoTurf":     "良馬場",
	"YayaomoTurf": "稍重馬場",
	"OmoTurf":     "重馬場",
	"HuryoTurf":   "不良馬場",
}

// ParseRequestInterface parseRequest関数を持つinterface
type ParseRequestInterface interface {
	parseRequest(channelSecret string, r events.APIGatewayProxyRequest) ([]*linebot.Event, error)
}

// ParseRequestStruct テストように作成
type ParseRequestStruct struct{}

func (p *ParseRequestStruct) parseRequest(channelSecret string, r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	event, err := ParseRequest(channelSecret, r)
	return event, err
}

// getLineMessage lineからのメッセージを取得する関数
func getLineMessage(p ParseRequestInterface, r events.APIGatewayProxyRequest) (lineMessage LineMessage) {
	lineMessage.Status = 200
	events, err := p.parseRequest(channelSecret, r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			lineMessage.Status = 400
		} else {
			lineMessage.Status = 500
		}
		fmt.Println(err)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			message := event.Message
			if value, ok := message.(*linebot.TextMessage); ok {
				lineMessage.Message = value
				lineMessage.ReplyToken = event.ReplyToken
			} else {
				lineMessage.Message = nil
			}
		}
	}

	return
}

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lineMessage := getLineMessage(&ParseRequestStruct{}, r)
	responseMessage := httpClientCourseResult(lineMessage.Message.Text)
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		lineMessage.Status = 500
	}
	if _, err = bot.ReplyMessage(lineMessage.ReplyToken, linebot.NewTextMessage(responseMessage)).Do(); err != nil {
		fmt.Printf("%v\n", err)
		lineMessage.Status = 500
	}
	return events.APIGatewayProxyResponse{
		Body:       r.Body,
		StatusCode: lineMessage.Status,
	}, nil
}

// httpClientGetId 馬名からidを取得する関数
// func httpClientGetId(horseName string) string {
// 	values, err := json.Marshal(HorseNameRequest{Category: "horse", HorseName: horseName})
// 	res, err := http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/horsename", "application/json", bytes.NewBuffer(values))
// 	if err != nil {
// 		fmt.Printf("%v\n", err)
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	horseNameData := &HorseNameData{}
// 	err = json.Unmarshal(body, horseNameData)
// 	if err != nil {
// 		fmt.Printf("failed unmarshal request: %v", err)
// 		log.Fatal(err)
// 	}

// 	return horseNameData.Data[0].ID
// }

// httpClientCourseResult idからその馬のコース成績を取得する関数
func httpClientCourseResult(name string) string {
	values, err := json.Marshal(CourseResultRequest{Name: name})
	res, err := http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/courseresult", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	CourseResultData := &CourseResultData{}
	err = json.Unmarshal(body, CourseResultData)
	if err != nil {
		fmt.Printf("failed unmarshal request: %v", err)
		log.Fatal(err)
	}
	v := reflect.ValueOf(CourseResultData.Data)
	responseMessage := ""

	distanceCount := 0
	babaCount := 0
	// レスポンスメッセージの生成
	for i := 0; i < v.NumField(); i++ {
		if _, ok := courseResultName[v.Type().Field(i).Name]; v.Field(i).Interface() != "0" && v.Type().Field(i).Name != "Name" && ok {
			responseMessage = responseMessage +
				courseResultName[v.Type().Field(i).Name] + ":" +
				v.Field(i).Interface().(string) + "\n"
		}

		if _, ok := distanceResultName[v.Type().Field(i).Name]; v.Field(i).Interface() != "0" && ok {
			if distanceCount == 0 {
				responseMessage = responseMessage + "(距離成績)\n"
			}

			responseMessage = responseMessage +
				distanceResultName[v.Type().Field(i).Name] + ":" +
				v.Field(i).Interface().(string) + "\n"

			distanceCount++
		}

		if _, ok := babaResultName[v.Type().Field(i).Name]; v.Field(i).Interface() != "0" && ok {
			if babaCount == 0 {
				responseMessage = responseMessage + "(馬場成績)\n"
			}

			responseMessage = responseMessage +
				babaResultName[v.Type().Field(i).Name] + ":" +
				v.Field(i).Interface().(string) + "\n"

			babaCount++
		}
	}

	return responseMessage
}

func createResposeMessage(resultName map[string]string, v reflect.Value, i int) (responseMessage string) {
	responseMessage = responseMessage +
		resultName[v.Type().Field(i).Name] + ":" +
		v.Field(i).Interface().(string) + "\n"

	return
}

func init() {
	util.EnvLoad()
	channelSecret = decryptKms(os.Getenv("CHANNEL_SECRET"))
	channelToken = decryptKms(os.Getenv("CHANNEL_TOKEN"))
}

func main() {
	lambda.Start(handler)
}
