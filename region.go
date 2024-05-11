package flowdata

import (
	"image/color"
)

const (
	divider float64 = 2
)

type Region struct {
	// name region
	Name string `json:"name"`
	// slice points
	Points [4]float64 `json:"points"`
	// region color
	Color color.Color `json:"color"`
}

// NewRegion init a new region
func NewRegion(name string, points [4]float64, color color.Color) *Region {
	return &Region{
		Name:   name,
		Points: points,
		Color:  color,
	}
}

// CheckEntranceInRegion checking an object for inclusion in a region
func (r *Region) CheckEntranceInRegion(objects []*Object) (in []*Object, out []*Object) {
	in = make([]*Object, 0, len(objects))
	out = make([]*Object, 0, len(objects))

	for key := range objects {
		result := r.bBoxInRegion(objects[key].BBox, r.Points)
		if result {
			in = append(in, objects[key])
		} else {
			out = append(out, objects[key])
		}
	}

	return
}

// bBoxInRegion accepts bbox in the format [x, y, w, h] and points in the format [x1, x2, y1, y2]
// returns true if bbox is inside the area and false otherwise
func (r *Region) bBoxInRegion(bbox [4]float64, points [4]float64) bool {
	centerX, centerY := CalculateBBoxCenter(bbox)

	return centerX >= points[0] && centerX <= points[1] && centerY >= points[2] && centerY <= points[3]
}

// calculateBBoxCenter accepts a bbox in the format [x, y, w, h]
// returns the coordinates of the center X, Y of the bbox
func CalculateBBoxCenter(bbox [4]float64) (centerX float64, centerY float64) {
	centerX = bbox[0] + bbox[2]/divider
	centerY = bbox[1] + bbox[3]/divider

	return
}
