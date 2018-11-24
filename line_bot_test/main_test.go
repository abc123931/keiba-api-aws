package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	mockhttp "github.com/karupanerura/go-mock-http-response"
)

func TestHandler(t *testing.T) {
	t.Run("getLineMessage正常終了", func(t *testing.T) {
		lineMessage := getLineMessage(&FakeParseRequest{}, events.APIGatewayProxyRequest{Body: `{"events":[{"type":"message","replyToken":"49c90bcfca564010aaf4b0800ce0238f","source":{"userId":"U7f4ea31bfb1249b50093074fc4cfd6b8","type":"user"},"timestamp":1540555234707,"message":{"type":"text","id":"8773115043446","text":"サトノダイヤモンド"}}]}`})
		expected := "サトノダイヤモンド"

		if !reflect.DeepEqual(lineMessage.Message.Text, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, lineMessage.Message.Text)
		}
	})

	t.Run("httpClientGetId正常終了", func(t *testing.T) {
		mockResponse(http.StatusOK, map[string]string{"Content-Type": "application/json"}, []byte(`{"data":[{"category":"horse","id":"13106101","name":"サトノダイヤモンド"}],"error":""}`))
		res := httpClientGetId("サトノダイヤモンド")
		expected := "13106101"
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})

	t.Run("httpClientCourseResult正常終了", func(t *testing.T) {
		mockResponse(http.StatusOK, map[string]string{"Content-Type": "application/json"}, []byte(`{"data":{"id":"13106101","sapporo_turf":"0","hakodate_turf":"0","fukushima_turf":"0","nigata_turf":"0","tokyo_turf":"0","nakayama_turf":"0","tyukyo_turf":"0","kyoto_turf":"(1-1-0-1)","hanshin_turf":"0","kokura_turf":"0","1000_turf":"(0-0-0-1)","1200_turf":"0","1400_turf":"0","1600_turf":"0","1800_turf":"0","2000_turf":"0","2200_turf":"0","2400_turf":"0","2500_turf":"0","3000_turf":"0","3200_turf":"0","3600_turf":"0"},"error":""}`))
		res := httpClientCourseResult("13106101")
		// fmt.Printf("courseResult: %v\n", res)
		expected := "京都成績:(1-1-0-1)\n\n" +
			"芝1000m:(0-0-0-1)\n"
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})
}

func mockResponse(statusCode int, headers map[string]string, body []byte) {
	http.DefaultClient = mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()
}
