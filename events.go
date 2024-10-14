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

	// Objects array
	Objects []*Object `json:"objects" gorm:"serializer:json;type:json"`
}

func (Event) TableName() string {
	return "events"
}
