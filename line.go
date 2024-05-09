package flowdata

type Line struct {
	Name   string    `json:"name"`
	Points []float64 `json:"points"`
	Color  string    `json:"color"`
}

func NewLine(name string, points []float64, color string) *Line {
	return &Line{name, points, color}
}

func (l *Line) CheckEntranceInRegion() {}
