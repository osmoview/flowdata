package flowdata

import "image/color"

type AnchorType int

const (
	Center AnchorType = iota
	Bottom
	Top
	AllPoints
)

type Line struct {
	// name region
	Name string `json:"name"`
	// array points
	A [2]float64
	// array points
	B [2]float64
	// region color
	Color color.Color `json:"color"`
}

// NewLine init a new line
func NewLine(name string, a, b [2]float64, color color.Color) *Line {
	return &Line{name, a, b, color}
}

// CheckObjectLineCrossing checks whether objects intersect a line and determines the direction of intersection
func (l *Line) CheckObjectLineCrossing(objects []*Object, anchor AnchorType) (crossedOneWay []*Object, crossedAnotherWay []*Object, notCrossed []*Object) {
	crossedOneWay = []*Object{}
	crossedAnotherWay = []*Object{}
	notCrossed = []*Object{}

	for _, object := range objects {
		var point [2]float64

		switch anchor {
		case Center:
			pointX, pointY := CalculateBBoxCenter(object.BBox)

			point = [2]float64{pointX, pointY}
		case Bottom:
			point = [2]float64{object.BBox[0] + object.BBox[2]/2, object.BBox[1] + object.BBox[3]}
		case Top:
			point = [2]float64{object.BBox[0] + object.BBox[2]/2, object.BBox[1]}
		case AllPoints:
			crossed := false
			directionSum := 0

			for _, p := range object.Area {
				if isCrossing, direction := l.isCrossingLine(p); isCrossing {
					crossed = true
					directionSum += direction
				}
			}

			if crossed {
				if directionSum > 0 {
					crossedOneWay = append(crossedOneWay, object)
				}
				if directionSum < 0 {
					crossedAnotherWay = append(crossedAnotherWay, object)
				}
			} else {
				notCrossed = append(notCrossed, object)
			}
		}

		if isCrossing, direction := l.isCrossingLine(point); isCrossing {
			if direction > 0 {
				crossedOneWay = append(crossedOneWay, object)
			}
			if direction < 0 {
				crossedAnotherWay = append(crossedAnotherWay, object)
			}
		} else {
			notCrossed = append(notCrossed, object)
		}
	}

	return
}

// isCrossingLine determines whether a point intersects a line
// returns the direction of intersection
func (l *Line) isCrossingLine(point [2]float64) (bool, int) {
	// Calculate the cross product to determine the position of a point relative to a line
	d := (point[0]-l.A[0])*(l.B[1]-l.A[1]) - (point[1]-l.A[1])*(l.B[0]-l.A[0])

	if d > 0 {
		// one way crossing
		return true, 1
	} else if d < 0 {
		// crossing in the other direction
		return true, -1
	}

	// does not cross
	return false, 0
}
