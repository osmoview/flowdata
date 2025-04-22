package flowdata

import (
	"time"
)

// Event struct contains the fields required for writing to the database
type Event struct {
	// event model contains ID as UUID, created and updated time
	ID        string    `json:"id" gorm:"primaryKey;unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// source UUID
	DataSource string `json:"datasource" gorm:"column:datasource"`

	// flow UUID
	Flow string `json:"flow"`

	// Status text message
	Status string `json:"status"`

	// Alert flag, if true then event should popup in the web interface and has some mark in the events archive
	Alert bool `json:"alert,omitempty"`

	// State of the event, like pass,fail,warn, or another
	State string `json:"state,omitempty"`

	// Objects array
	Objects []*Object `json:"objects" gorm:"serializer:json;type:json"`

	// regions and lines used to create this event
	Regions []*Region `json:"regions" gorm:"serializer:json;type:json"`
	Lines   []*Line   `json:"lines" gorm:"serializer:json;type:json"`

	// data may contains any required custom data, like some variable etc
	Data map[string]any `json:"data" gorm:"serializer:json;type:json"`
}

func (Event) TableName() string {
	return "events"
}
