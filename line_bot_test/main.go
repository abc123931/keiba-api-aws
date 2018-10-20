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
	Category string `json:"category"`
	Name     string `json:"horse_name"`
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
	httpClient()
	lineMessage := getLineMessage(r)
	fmt.Printf("%v", lineMessage)
	fmt.Printf("%v", lineMessage.Message.Text)
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = bot.ReplyMessage(lineMessage.ReplyToken, linebot.NewTextMessage(lineMessage.Message.Text)).Do(); err != nil {
		log.Print(err)
	}
	return events.APIGatewayProxyResponse{
		Body:       r.Body,
		StatusCode: lineMessage.Status,
	}, nil
}

func httpClient() {
	values, err := json.Marshal(HorseNameRequest{Category: "horse", Name: "サトノダイヤモンド"})
	res, err := http.Post("https://xs8k217r0j.execute-api.ap-northeast-1.amazonaws.com/Prod/horsename", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", body)
	fmt.Printf("%v", string(body))
}

func init() {
	util.EnvLoad()
	channelSecret = decryptKms(os.Getenv("CHANNEL_SECRET"))
	channelToken = decryptKms(os.Getenv("CHANNEL_TOKEN"))
}

func main() {
	lambda.Start(handler)
}
