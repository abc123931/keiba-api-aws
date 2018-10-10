package main

import (
	"fmt"
	"log"
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

func init() {
	util.EnvLoad()
	channelSecret = decryptKms(os.Getenv("CHANNEL_SECRET"))
	channelToken = decryptKms(os.Getenv("CHANNEL_TOKEN"))
}

func main() {
	lambda.Start(handler)
}
