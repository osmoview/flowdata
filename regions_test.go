package flowdata

import "testing"

func TestRegionBBox(t *testing.T) {
	// triangle
	r := Region{
		Points: [][2]float32{
			{0.1, 0.9},
			{0.5, 0.1},
			{0.9, 0.8},
		},
	}

	bbox, ok := r.BBox()
	if !ok {
		t.Error("no bbox")
	}

	if bbox[0] != 0.1 || bbox[1] != 0.1 || bbox[2] != 0.8 || bbox[3] != 0.8 {
		t.Errorf("invalid bbox: %v", bbox)
	}
}

func TestRegionContainsPoint(t *testing.T) {
	r := Region{
		Points: [][2]float32{
			{0.1, 0.9},
			{0.5, 0.1},
			{0.9, 0.8},
			{0.5, 0.51},
		},
	}

	if !r.ContainsPoint([2]float32{0.5, 0.5}) {
		t.Error("point in the center outside region")
	}
	if p := [2]float32{0.2, 0.2}; r.ContainsPoint(p) {
		t.Errorf("%v point inside region", p)
	}
	if p := [2]float32{0.7, 0.3}; r.ContainsPoint(p) {
		t.Errorf("%v point inside region", p)
	}
	if p := [2]float32{0.5, 0.52}; r.ContainsPoint(p) {
		t.Errorf("%v point inside region", p)
	}
}
