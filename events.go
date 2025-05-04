package main

type UserEvent struct {
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"timestamp"`
}
