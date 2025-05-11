package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type UserEvent struct {
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"timestamp"`
}

func getTLSConfig() *tls.Config {
	caCert, err := os.ReadFile("secrets/ca-cert.pem")
	if err != nil {
		log.Fatal("Failed to read CA certificate:", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	return &tls.Config{RootCAs: certPool}
}

func main() {
	log.Println("Starting auth service...")
	dialer := &kafka.Dialer{
		TLS: getTLSConfig(),
		SASLMechanism: plain.Mechanism{
			Username: "testuser",
			Password: "testpassword",
		},
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9093"},
		Topic:   "user-events",
		Dialer:  dialer,
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
