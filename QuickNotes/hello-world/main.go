package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	db = dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})))
	tableName = "NotesTable"
)

type Note struct{
	NoteID string `json:"noteId"`
	Content string `json:"content"`
}
func handler(ctx context.Context,request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
			return getNoteHandler(request)
	default:
	return events.APIGatewayProxyResponse{
		Body:       "Method not allowed",
		StatusCode: 200,
	}, nil
	}
}

func getNoteHandler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
	noteID := request.QueryStringParameters["noteId"]
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(noteID),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Failed to get note: %s",err.Error()),
		}, nil
	}
	if result.Item == nil{
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body: fmt.Sprintf("Note not found"),
		}, nil
	}

	var note Note
	err = dynamodbattribute.UnmarshalMap(result.Item,&note)
	if err != nil{
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Failed to Unmarshal note: %s",err.Error()),
		}, nil
	}
	body, err := json.Marshal(note)
	if err != nil{
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Failed to marshal note: %s",err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:  string (body),
	}, nil
}



func main() {
	lambda.Start(handler)
}
