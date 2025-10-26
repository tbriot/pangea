package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/go-github/v76/github"
)

var (
	secretToken = []byte(os.Getenv("GITHUB_WEBHOOK_SECRET"))
	sqsClient   *sqs.Client
	queueURL    string
)

func init() {
	// Initialize AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, ERROR: %v", err)
	}

	// Initialize SQS client
	sqsClient = sqs.NewFromConfig(cfg)

	// Get queue URL from environment variable
	queueURL = os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, secretToken)
	if err != nil {
		fmt.Printf("GitHub Webhook event request is invalid: %v\n", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		fmt.Printf("Error while parsing the GitHub webhook event request payload: %v", err)
		http.Error(w, "Unable to parse webhook event payload", http.StatusBadRequest)
		return
	}
	switch event := event.(type) {
	case *github.PushEvent:
		err := pushEventToSqsQueue(event)
		if err != nil {
			fmt.Printf("Error processing push event: %v\n", err)
			http.Error(w, "Failed to process push event", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Push event processed successfully")
	default:
		fmt.Printf("Event %T ignored\n", event)
		http.Error(w, "Only 'push' event are authorized. Discarding event.", http.StatusUnauthorized)
		return
	}
}

func pushEventToSqsQueue(e *github.PushEvent) error {
	// Convert the GitHub push event to JSON
	eventJSON, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal push event: %w", err)
	}

	// Create the SQS message input
	input := &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(eventJSON)),
	}

	// Send the message to SQS
	result, err := sqsClient.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	fmt.Printf("Successfully sent push event to SQS. Message ID: %s\n", *result.MessageId)
	return nil
}

func main() {
	http.HandleFunc("/event", eventHandler)
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
