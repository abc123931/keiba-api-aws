package util

import (
	"fmt"

	"github.com/joho/godotenv"
)

// EnvLoad .envを読み込む関数
func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("start lambda function at prod")
	}
}
