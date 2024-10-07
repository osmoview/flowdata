package flowdata

import "github.com/google/uuid"

type Region struct {
	ID     string       `json:"id" yaml:"id"`
	Name   string       `json:"name" yaml:"name"`
	Points [][2]float64 `json:"points" yaml:"points,flow"`
	Color  string       `json:"color,omitempty" yaml:"color,omitempty"`
}

// NewRegion creates new region instance
func NewRegion(name string, points [][2]float64, color string) *Region {
	return &Region{
		ID:     uuid.NewString(),
		Name:   name,
		Points: points,
		Color:  color,
	}
}

// BBox transform regions points to the bbox[x,y,w,h] format
func (r *Region) BBox() (BBox, bool) {
	if len(r.Points) < 3 {
		return [4]float64{0, 0, 0, 0}, false
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

	return [4]float64{minX, minY, width, height}, true
}

// ContainsBBox checks if region contains all bbox points
func (r *Region) ContainsBBox(bbox BBox) bool {
	// check top left point
	if !r.ContainsPoint([2]float64{bbox[0], bbox[1]}) {
		return false
	}

	// check top right point
	if !r.ContainsPoint([2]float64{bbox[0] + bbox[2], bbox[1]}) {
		return false
	}

	// check bottom right point
	if !r.ContainsPoint([2]float64{bbox[0] + bbox[2], bbox[1] + bbox[3]}) {
		return false
	}

	// check bottom left point
	if !r.ContainsPoint([2]float64{bbox[0], bbox[1] + bbox[3]}) {
		return false
	}

	return true
}

// ContainsPoint checks if point inside a region
func (r *Region) ContainsPoint(p [2]float64) bool {
	// create a new point outside the area
	out := [2]float64{p[0] + 1, (p[1] + 1) * 1.1}

	// raycast algorithm: check intersections with each lines from region to line with line p-out
	// if number of intersections is odd, then point inside region
	var intersections = 0
	for i := 0; i < len(r.Points); i++ {
		r1 := r.Points[i]

		j := i + 1
		if j >= len(r.Points) {
			j = 0
		}

		// r1-r2 is line in the region
		r2 := r.Points[j]

		// get intersection point
		x, ok := GetIntersectionPointOfTwoLines(r1, r2, p, out)
		if ok && !(x[0] == r1[0] && x[1] == r2[1]) {
			intersections++
		}
	}

	return intersections%2 != 0
}

func GetIntersectionPointOfTwoLines(a1, a2, b1, b2 [2]float64) ([2]float64, bool) {
	den := (b2[1]-b1[1])*(a2[0]-a1[0]) - (b2[0]-b1[0])*(a2[1]-a1[1])

	// lines are parallel
	if den == 0 {
		return [2]float64{}, false
	}

	ua := ((b2[0]-b1[0])*(a1[1]-b1[1]) - (b2[1]-b1[1])*(a1[0]-b1[0])) / den
	ub := ((a2[0]-a1[0])*(a1[1]-b1[1]) - (a2[1]-a1[1])*(a1[0]-b1[0])) / den

	// is the intersection along the segments
	if ua <= 0 || ua > 1 || ub <= 0 || ub > 1 {
		return [2]float64{}, false
	}

	// return a object with the x and y coordinates of the intersection
	return [2]float64{a1[0] + ua*(a2[0]-a1[0]), a1[1] + ua*(a2[1]-a1[1])}, true
}
