package flowdata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	JPEGExt = ".jpg"
	JSONExt = ".json"
)

var EventsDir = "./data/events"

// Event struct contains the fields required for writing to the database
type Event struct {
	// event model contains ID as UUID, created and updated time
	ID        string    `json:"id" gorm:"primaryKey;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index"`

	// source UUID
	DataSource string `json:"datasource" gorm:"index;column:datasource"`

	// flow UUID
	Flow string `json:"flow" gorm:"index"`

	// Status text message
	Status string `json:"status" gorm:"index"`

	// Alert flag, if true then event should popup in the web interface and has some mark in the events archive
	Alert bool `json:"alert,omitempty" gorm:"index"`

	// State of the event, like pass,fail,warn, or another
	State string `json:"state,omitempty" gorm:"index"`

	// events group name
	Group string `json:"group,omitempty" gorm:"index"`

	// Objects array
	Objects []*Object `json:"objects" gorm:"serializer:json;type:json"`

	// regions and lines used to create this event
	Regions []*Region `json:"regions,omitempty" gorm:"serializer:json;type:json"`
	Lines   []*Line   `json:"lines,omitempty" gorm:"serializer:json;type:json"`

	// data may contains any required custom data, like some variable etc
	Data map[string]any `json:"data,omitempty" gorm:"serializer:json;type:json"`

	// archive_file contains path to video file in the video archive
	ArchiveFile string `json:"archive_file,omitempty" gorm:"-"`

	Image string `json:"image,omitempty" gorm:"image"`
}

func (Event) TableName() string {
	return "events"
}

func (e *Event) Save(frame *Frame) error {
	abs, rel := e.buildFilename()

	// save frame image to disk
	// TODO: check frame extension
	if frame != nil {
		e.Image = rel + JPEGExt
		if err := e.SaveFrame(abs+JPEGExt, frame); err != nil {
			return fmt.Errorf("failed to save frame image: %v", err)
		}
	}

	// if objects exists, write it to disk in the JSON format
	if len(e.Objects) > 0 {
		if err := e.SaveObjects(abs+JSONExt, e.Objects); err != nil {
			log.Printf("failed to save objects: %v", err)
		}
	}

	return nil
}

func (e *Event) SaveFrame(fn string, frame *Frame) error {
	imgdata, err := frame.GetBytes()
	if err != nil {
		return err
	}

	return os.WriteFile(fn, imgdata, 0666)
}

func (e *Event) SaveObjects(fn string, objects []*Object) error {
	b, err := json.Marshal(objects)
	if err != nil {
		return err
	}

	return os.WriteFile(fn, b, 0666)
}

func (e *Event) getDir() string {
	if dir := os.Getenv("OSMO_EVENTS_DIR"); dir != "" {
		return dir
	}

	return EventsDir
}

func (e *Event) buildFilename() (abs string, rel string) {
	dir := e.getDir()

	// relative path to the directory where save files
	reldir := filepath.Join(e.DataSource, e.Flow)

	// nolint:errcheck
	os.MkdirAll(filepath.Join(dir, reldir), 0755)

	// absolute and relative path to the file, without extension
	abs = filepath.Join(dir, reldir, e.ID)
	rel = filepath.Join(reldir, e.ID)

	return
}
