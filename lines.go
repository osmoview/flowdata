package flowdata

import "github.com/google/uuid"

const (
	DirectionBoth     = "both"
	DirectionForward  = "forward"
	DirectionBackward = "backward"
)

type Line struct {
	ID        string       `json:"id" yaml:"id"`
	Name      string       `json:"name" yaml:"name"`
	Points    [][2]float64 `json:"points" yaml:"points"`
	Color     string       `json:"color,omitempty" yaml:"color,omitempty"`
	Direction string       `json:"direction"`
}

func NewLine(name string, points [][2]float64, color string) *Line {
	return &Line{
		ID:     uuid.NewString(),
		Name:   name,
		Points: points,
		Color:  color,
	}
}
