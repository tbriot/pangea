package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v76/github"
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		fmt.Println("Payload validation error: ", err)
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		fmt.Println("Event parsing error:", err)
	}
	switch event := event.(type) {
	case *github.PushEvent:
		processPushEvent(event)
	default:
		fmt.Printf("Event %T ignored\n", event)

		w.WriteHeader(http.StatusOK)
	}
}

func processPushEvent(e *github.PushEvent) {
	fmt.Println("Processing push event")
	after := e.GetAfter()
	pusher := e.GetPusher().String()
	fmt.Println("SHA of the most recent commit=", after)
	fmt.Println("Git commiter=", pusher)
}

func main() {
	http.HandleFunc("/", eventHandler)
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
