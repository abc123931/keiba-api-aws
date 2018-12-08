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

	t.Run("httpClientCourseResult正常終了", func(t *testing.T) {
		mockResponse(http.StatusOK, map[string]string{"Content-Type": "application/json"}, []byte(`{"data":{"name":"ブラストワンピース","sapporo_turf":"0","hakodate_turf":"0","hukushima_turf":"0","nigata_turf":"0","tokyo_turf":"0","nakayama_turf":"0","tyukyo_turf":"0","kyoto_turf":"(1-1-0-1)","hanshin_turf":"0","kokura_turf":"0","1000_turf":"(0-0-0-1)","1200_turf":"0","1400_turf":"0","1600_turf":"0","1800_turf":"0","2000_turf":"0","2200_turf":"0","2400_turf":"0","2500_turf":"0","3000_turf":"0","3200_turf":"0","3600_turf":"0","ryo_turf":"(1-1-1-1)","yayaomo_turf":"0","omo_turf":"0","huryo_turf":"0"},"error":""}`))
		res := httpClientCourseResult("13106101")
		expected := "(コース成績)\n京都成績:(1-1-0-1)\n" +
			"\n(距離成績)\n芝1000m:(0-0-0-1)\n" +
			"\n(馬場成績)\n良馬場:(1-1-1-1)\n"
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})

	t.Run("requestRaceIndexApi正常終了", func(t *testing.T) {
		mockResponse(http.StatusOK,
			map[string]string{"Content-Type": "application/json"},
			[]byte(`{"data":[{"horse_name":"アーモンドアイ",`+
				`"total_index":"83.4","train_index":"21.3","stable_index":"12.9"},`+
				`{"horse_name":"サトノダイヤモンド","total_index":"78.4",`+
				`"train_index":"21.3","stable_index":"6.9"}],"error":""}`))

		res := requestRaceIndexAPI("ジャパンカップ")
		expected := []RaceIndex{}
		expected = append(expected,
			RaceIndex{Name: "アーモンドアイ", TotalIndex: "83.4",
				TrainIndex: "21.3", StableIndex: "12.9"})

		expected = append(expected,
			RaceIndex{Name: "サトノダイヤモンド", TotalIndex: "78.4",
				TrainIndex: "21.3", StableIndex: "6.9"})

		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})

	t.Run("requestRaceIndexApiレスポンス空", func(t *testing.T) {
		mockResponse(http.StatusOK,
			map[string]string{"Content-Type": "application/json"},
			[]byte(`{"data":[],"error":""}`))

		res := requestRaceIndexAPI("ジャパンカップ")
		expected := []RaceIndex{}
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("response not same expected: expected is %v response is  %v", expected, res)
		}
	})
}

func mockResponse(statusCode int, headers map[string]string, body []byte) {
	http.DefaultClient = mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()
}
