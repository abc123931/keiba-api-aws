package main

import (
	"fmt"
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
	Message *linebot.TextMessage
	Status  int
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
	return events.APIGatewayProxyResponse{
		Body:       r.Body,
		StatusCode: lineMessage.Status,
	}, nil
}

func init() {
	util.EnvLoad()
	channelSecret = os.Getenv("CHANNEL_SECRET")
	channelToken = os.Getenv("CHANNEL_TOKEN")
}

func main() {
	lambda.Start(handler)
}
