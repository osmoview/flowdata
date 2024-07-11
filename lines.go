package flowdata

type Line struct {
	Name   string       `json:"name" yaml:"name"`
	Points [][2]float64 `json:"points" yaml:"points"`
	Color  string       `json:"color,omitempty" yaml:"color,omitempty"`
}

func NewLine(name string, points [][2]float64, color string) *Line {
	return &Line{
		Name:   name,
		Points: points,
		Color:  color,
	}
}
