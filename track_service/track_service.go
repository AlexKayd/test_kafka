package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"os"

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
	dialer := &kafka.Dialer{
		TLS: getTLSConfig(),
		SASLMechanism: plain.Mechanism{
			Username: "testuser",
			Password: "testpassword",
		},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9093"},
		Topic:   "user-events",
		GroupID: "track-service-group",
		Dialer:  dialer,
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
