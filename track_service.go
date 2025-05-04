package main

import (
	"context"
	"encoding/json"
	"log"
	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-events",
		GroupID: "track-service-group",
	})
	defer reader.Close()

	log.Println("Track service listening for user events...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("error reading message:", err)
		}

		var event UserEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("failed to decode message:", err)
			continue
		}

		if event.EventType == "login" {
			log.Printf("User %s logged in.\n", event.UserID)
		}
	}
}
