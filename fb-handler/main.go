package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FbWebhookEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Time    int64  `json:"time"`
		Changes []struct {
			Field string `json:"field"`
			Value struct {
				From struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"from"`
				Post struct {
					StatusType      string `json:"status_type"`
					IsPublished     bool   `json:"is_published"`
					UpdatedTime     string `json:"updated_time"`
					PermalinkURL    string `json:"permalink_url"`
					PromotionStatus string `json:"promotion_status"`
					ID              string `json:"id"`
				} `json:"post"`
				Message     string `json:"message"`
				PostID      string `json:"post_id"`
				CommentID   string `json:"comment_id"`
				CreatedTime int64  `json:"created_time"`
				Item        string `json:"item"`
				ParentID    string `json:"parent_id"`
				Verb        string `json:"verb"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event FbWebhookEvent

	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid request body: %v"}`, err),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	messages := []string{}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.UTC
	}

	for _, entry := range event.Entry {
		for _, change := range entry.Changes {
			t := time.Unix(entry.Time, 0).UTC().In(location)
			timeStr := t.Format(time.RFC3339)
			msg := fmt.Sprintf("%s: %s : %s", change.Value.From.Name, change.Value.Message, timeStr)
			messages = append(messages, msg)
		}
	}

	if len(messages) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"message":"No comments found"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	respBody, _ := json.Marshal(map[string]interface{}{
		"messages": messages,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(respBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
