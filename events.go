package flowdata

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var TableName = "events"

// Event struct contains the fields required for writing to the database
type Event struct {
	// event model contains ID as UUID, created and updated time
	ID        string    `json:"id" gorm:"primaryKey;unique;type:uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// source UUID
	DataSource string `json:"datasource"`

	// flow UUID
	Flow string `json:"flow"`

	// Status text message
	Status string `json:"status"`

	// Objects array
	Objects Objects `json:"objects" gorm:"type:json"`
}

func (Event) TableName() string {
	return TableName
}

//
// objects custom data type implements scanner and valuer methods
//

type Objects []*Object

func (o Objects) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *Objects) Scan(val any) error {
	s, ok := val.(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value")
	}

	return json.Unmarshal([]byte(s), o)
}
