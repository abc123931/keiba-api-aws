package util

import (
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("ConnectTable正常終了", func(t *testing.T) {
		res := ConnectTable("test_table", "testRegion", "testEndPoint")

		if reflect.TypeOf(res).String() != "dynamo.Table" {
			t.Fatalf("response not dynamo.Table:  %T", res)
		}
	})

}
