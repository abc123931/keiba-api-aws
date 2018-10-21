package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	ID string `json:"id"`
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
	ID                    string `dynamo:"id" json:"id"`
	SapporoTurf           string `dynamo:"sapporo_turf" json:"sapporo_turf"`
	HakodateTurf          string `dynamo:"hakodate_turf" json:"hakodate_turf"`
	FukushimaTurf         string `dynamo:"fukushima_turf" json:"fukushima_turf"`
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
}

// getLineMessage lineからのメッセージを取得する関数
func getLineMessage(r events.APIGatewayProxyRequest) (lineMessage LineMessage) {
	lineMessage.Status = 200
	events, err := ParseRequest(channelSecret, r)
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
	lineMessage := getLineMessage(r)
	responseMessage := httpClient(lineMessage.Message.Text)
	fmt.Printf("%v", lineMessage)
	fmt.Printf("%v", lineMessage.Message.Text)
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = bot.ReplyMessage(lineMessage.ReplyToken, linebot.NewTextMessage(responseMessage)).Do(); err != nil {
		log.Print(err)
	}
	return events.APIGatewayProxyResponse{
		Body:       r.Body,
		StatusCode: lineMessage.Status,
	}, nil
}

func httpClient(horseName string) string {
	values, err := json.Marshal(HorseNameRequest{Category: "horse", HorseName: horseName})
	res, err := http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/horsename", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	horseNameData := &HorseNameData{}
	err = json.Unmarshal(body, horseNameData)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed unmarshal request: %v", err)
	}

	values, err = json.Marshal(CourseResultRequest{ID: horseNameData.Data[0].ID})
	res, err = http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/courseresult", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	CourseResultData := &CourseResultData{}
	err = json.Unmarshal(body, CourseResultData)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed unmarshal request: %v", err)
	}

	responseMessage := "札幌成績:" + CourseResultData.Data.SapporoTurf + "\n" +
		"函館成績:" + CourseResultData.Data.HakodateTurf + "\n" +
		"福島成績:" + CourseResultData.Data.FukushimaTurf + "\n" +
		"新潟成績:" + CourseResultData.Data.NigataTurf + "\n" +
		"東京成績:" + CourseResultData.Data.TokyoTurf + "\n" +
		"中山成績:" + CourseResultData.Data.NakayamaTurf + "\n" +
		"中京成績:" + CourseResultData.Data.TyukyoTurf + "\n" +
		"京都成績:" + CourseResultData.Data.KyotoTurf + "\n" +
		"阪神成績:" + CourseResultData.Data.HanshinTurf + "\n" +
		"小倉成績:" + CourseResultData.Data.KokuraTurf + "\n\n" +
		"(距離成績)" + "\n" +
		"芝1000m:" + CourseResultData.Data.ThousandTurf + "\n" +
		"芝1200m:" + CourseResultData.Data.TwelveHundredTurf + "\n" +
		"芝1400m:" + CourseResultData.Data.FourteenHundredTurf + "\n" +
		"芝1600m:" + CourseResultData.Data.SixteenHundredTurf + "\n" +
		"芝1800m:" + CourseResultData.Data.EighteenHundredTurf + "\n" +
		"芝2000m:" + CourseResultData.Data.TwoThousandTurf + "\n" +
		"芝2200m:" + CourseResultData.Data.TwentyTwoHundredTurf + "\n" +
		"芝2400m:" + CourseResultData.Data.TwentyFourHundredTurf + "\n" +
		"芝2500m:" + CourseResultData.Data.TwentyFiveHundredTurf + "\n" +
		"芝3000m:" + CourseResultData.Data.ThreeThousandTurf + "\n" +
		"芝3200m:" + CourseResultData.Data.ThirtyTwoHundredTurf + "\n" +
		"芝3600m:" + CourseResultData.Data.ThirtySixHundredTurf

	return responseMessage
}

func init() {
	util.EnvLoad()
	channelSecret = decryptKms(os.Getenv("CHANNEL_SECRET"))
	channelToken = decryptKms(os.Getenv("CHANNEL_TOKEN"))
}

func main() {
	lambda.Start(handler)
}
