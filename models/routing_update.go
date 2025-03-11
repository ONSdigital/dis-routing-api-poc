package models

import "time"

type RoutingUpdatedEvent struct {
	EventType string    `json:"event_type"` // e.g., "route_updates"
	Changes   []Changes `json:"changes"`    // List of batched events
}

// Changes represents an individual event within the batch
type Changes struct {
	Action    string      `json:"action"`    // "create", "update", "delete"
	Entity    string      `json:"entity"`    // "route", "redirect"
	ID        string      `json:"id"`        // Document ID
	Payload   interface{} `json:"payload"`   // The actual update
	Timestamp time.Time   `json:"timestamp"` // When the event was generated
}
