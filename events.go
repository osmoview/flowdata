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
	var data []byte

	switch v := val.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("%T is unexpected data type", val)
	}

	return json.Unmarshal(data, o)
}
