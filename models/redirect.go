package models

import "time"

// Redirect represents a path-based redirect.
type Redirect struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	From       string    `bson:"from" json:"from"`
	To         string    `bson:"to" json:"to"`
	StatusCode int       `bson:"status_code" json:"status_code"` // 307 or 308 only
	Domains    []string  `bson:"domains" json:"domains"`
	Subnet     []string  `bson:"subnet" json:"subnet"`
	ReviewDate time.Time `bson:"review_date" json:"review_date"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
