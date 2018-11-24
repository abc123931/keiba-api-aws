package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

// FakeParseRequest テスト用の構造体
type FakeParseRequest struct{}

func (p *FakeParseRequest) parseRequest(channelSecret string, r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), request); err != nil {
		return nil, err
	}
	return request.Events, nil
}
