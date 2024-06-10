package flowdata

import "math"

type Region struct {
	Name   string       `json:"name,omitempty"`
	Points [][2]float64 `json:"points"`
	Color  string       `json:"color,omitempty"`
}

func NewRegion(name string, points [][2]float64, color string) *Region {
	return &Region{
		Name:   name,
		Points: points,
		Color:  color,
	}
}

// BBox get a rectangle by points
func (r *Region) BBox() (minX, minY, maxX, maxY float64) {
	if len(r.Points) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY = math.MaxFloat64, math.MaxFloat64
	maxX, maxY = math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64

	for _, point := range r.Points {
		if point[0] < minX {
			minX = point[0]
		}
		if point[0] > maxX {
			maxX = point[0]
		}
		if point[1] < minY {
			minY = point[1]
		}
		if point[1] > maxY {
			maxY = point[1]
		}
	}

	return minX, minY, maxX, maxY
}
