package flowdata

import (
	"testing"
)

func TestRegion_BBox(t *testing.T) {
	tests := []struct {
		name     string
		points   [][2]float64
		expected [4]float64
	}{
		{
			name:     "empty points",
			points:   [][2]float64{},
			expected: [4]float64{0, 0, 0, 0},
		},
		{
			name:     "single point",
			points:   [][2]float64{{10, 20}},
			expected: [4]float64{1, 1, 9, 19},
		},
		{
			name:     "multiple points",
			points:   [][2]float64{{10, 20}, {30, 40}, {5, 15}},
			expected: [4]float64{1, 1, 29, 39},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegion("", tt.points, "")
			bbox := r.BBox()

			if bbox[0] != tt.expected[0] || bbox[1] != tt.expected[1] || bbox[2] != tt.expected[2] || bbox[3] != tt.expected[3] {
				t.Errorf(
					"expected bounding box (%f, %f, %f, %f), got (%f, %f, %f, %f)",
					tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3],
					bbox[0], bbox[1], bbox[2], bbox[3],
				)
			}
		})
	}
}
