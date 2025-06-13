package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sourceIP := request.RequestContext.Identity.SourceIP
	if sourceIP == "" {
		sourceIP = "world"
	}

	// สร้าง JSON response
	responseBody, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("Hello, %s!", sourceIP),
	})

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
