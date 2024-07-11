package flowdata

type Region struct {
	Name   string       `json:"name" yaml:"name"`
	Points [][2]float64 `json:"points" yaml:"points"`
	Color  string       `json:"color,omitempty" yaml:"color,omitempty"`
}

func NewRegion(name string, points [][2]float64, color string) *Region {
	return &Region{
		Name:   name,
		Points: points,
		Color:  color,
	}
}

// BBox transform regions points to the bbox[x,y,w,h] format
func (r *Region) BBox() [4]float64 {
	if len(r.Points) == 0 {
		return [4]float64{0, 0, 0, 0}
	}

	minX, minY := 1.0, 1.0
	maxX, maxY := 0.0, 0.0

	for _, point := range r.Points {
		if point[0] < minX {
			minX = point[0]
		}
		if point[1] < minY {
			minY = point[1]
		}
		if point[0] > maxX {
			maxX = point[0]
		}
		if point[1] > maxY {
			maxY = point[1]
		}
	}

	width := maxX - minX
	height := maxY - minY

	return [4]float64{minX, minY, width, height}
}
