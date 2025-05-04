package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-events",
	})
	defer writer.Close()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user_id", http.StatusBadRequest)
			return
		}

		event := UserEvent{
			UserID:    userID,
			EventType: "login",
			Timestamp: time.Now().Unix(),
		}
		value, _ := json.Marshal(event)

		err := writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(userID),
			Value: value,
		})
		if err != nil {
			log.Println("failed to write message:", err)
			http.Error(w, "kafka error", 500)
			return
		}

		log.Println("Sent login event for user:", userID)
		w.Write([]byte("Login successful"))
	})

	log.Println("Auth service listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
