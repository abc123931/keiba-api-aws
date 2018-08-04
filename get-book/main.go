package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
)

var (
	dynamoRegion   string
	dynamoEndpoint string
)

type Book struct {
	Title    string `dynamo:"title" json:"title"`
	Category string `dynamo:"category" json:"category"`
}

type DbConnect interface {
	get(title string, category string, book *Book) (err error)
}

type Table struct {
	Table dynamo.Table
}

type Request struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("start lambda function at prod")
	}
}

func (table *Table) get(title string, category string, book *Book) (err error) {
	err = table.Table.Get("title", title).
		Range("category", dynamo.Equal, category).
		One(book)
	return
}

func getBook(db DbConnect, title string, category string) (response string, err error) {
	book := Book{}
	err = db.get(title, category, &book)
	if err != nil {
		log.Println(err)
		return
	}

	json, err := json.Marshal(book)
	if err != nil {
		log.Println(err)
		return
	}
	response = string(json)

	return
}

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String(dynamoRegion),
		Endpoint: aws.String(dynamoEndpoint),
	})

	table := &Table{db.Table("Books")}

	request := &Request{}
	err := json.Unmarshal([]byte(r.Body), request)

	response := ""

	if err != nil {
		response = fmt.Sprintf("cannot encode request json: %s", err)
	} else {
		response, err = getBook(table, request.Title, request.Category)
		if err != nil {
			response = fmt.Sprintf("cannot get book: %s", err)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       response,
		StatusCode: 200,
	}, nil
}

func init() {
	Env_load()
	dynamoRegion = os.Getenv("DYNAMO_REGION")
	dynamoEndpoint = os.Getenv("DYNAMO_ENPOINT")
}
func main() {
	lambda.Start(handler)
}
