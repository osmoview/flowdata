package flowdata

import (
	"errors"
	"math"
)

var ErrMismatchVectorsLength = errors.New("vectors lengths mismatch")
var ErrVectorNegativeValue = errors.New("l2dist has negative value")

type Vector []float32

// Empty check
func (v Vector) IsValid() bool {
	return len(v) > 0
}

// EuclideanDistance calculate the distance between two vectors
// dist = sqrt((x2-x1)^2 + (y2-y1)^2 + ... + (n2-n1)^2)
func (v Vector) EuclideanDistance(b Vector) (float32, error) {
	if len(v) != len(b) {
		return -1, ErrMismatchVectorsLength
	}

	var sum float32
	for i := 0; i < len(v); i++ {
		x := b[i] - v[i]
		sum += x * x
	}

	dist := float32(math.Sqrt(float64(sum)))
	if dist < 0 {
		return dist, ErrVectorNegativeValue
	}

	return dist, nil
}
