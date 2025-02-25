package models

import "time"

// Route represents a path-based route.
type Route struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Path        string    `bson:"path" json:"path"`
	Method      string    `bson:"method" json:"method"`
	Destination string    `bson:"destination" json:"destination"`
	Domains     []string  `bson:"domains" json:"domains"`
	Subnet      []string  `bson:"subnet" json:"subnet"` // web, publishing, or both
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
