package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ParseRequest(channelSecret string, r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	fmt.Printf("%v\n", r.Headers)
	fmt.Printf("%v\n", r.Headers["X-Line-Signature"])
	fmt.Printf("%v\n", r.Body)
	fmt.Printf("%v\n", channelSecret)
	kmsClient := kms.New(session.New())
	decodedBytes, err := base64.StdEncoding.DecodeString(channelSecret)
	if err != nil {
		panic(err)
	}
	input := &kms.DecryptInput{
		CiphertextBlob: decodedBytes,
	}
	response, err := kmsClient.Decrypt(input)
	if err != nil {
		panic(err)
	}
	// Plaintext is a byte array, so convert to string
	decrypted = string(response.Plaintext[:])
	fmt.Printf("%v\n", decrypted)
	if !validateSignature(decrypted, r.Headers["X-Line-Signature"], []byte(r.Body)) {
		fmt.Println("シグネチャの検証に失敗しました。")
		return nil, linebot.ErrInvalidSignature
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}

	if err := json.Unmarshal([]byte(r.Body), request); err != nil {
		return nil, err
	}
	return request.Events, nil
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	hash.Write(body)
	return hmac.Equal(decoded, hash.Sum(nil))
}
