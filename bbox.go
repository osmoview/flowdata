package flowdata

import "math"

// BBox describe position and size of the object [x,y,w,h] in relative to the image coordinates
type BBox [4]float32

// IsValid checks whether the bbox has width and heigth
func (bbox BBox) IsValid() bool {
	return bbox[2] > 0 && bbox[3] > 0
}

// Scale bbox by specified value
func (bbox BBox) Scale(n float32) BBox {
	kw := bbox[2]*n - bbox[2]
	kh := bbox[3]*n - bbox[3]

	b := BBox{
		bbox[0] - kw/2,
		bbox[1] - kh/2,
		bbox[2] + kw,
		bbox[3] + kh,
	}

	return b.Fix()
}

// Fix should fit the bbox in the 0.0-1.0 range
func (bbox BBox) Fix() (b BBox) {
	b = bbox
	if b[0] < 0 {
		b[0] = 0
	}
	if b[1] < 0 {
		b[1] = 0
	}
	if b[0]+b[2] > 1 {
		b[2] = 1 - b[0]
	}
	if b[1]+b[3] > 1 {
		b[3] = 1 - b[1]
	}

	return b
}

// Center returns center point from the bbox
func (bbox BBox) Center() [2]float32 {
	return [2]float32{
		bbox[0] + bbox[2]/2,
		bbox[1] + bbox[3]/2,
	}
}

// Bottom returns center bottom point
func (bbox BBox) Bottom() [2]float32 {
	return [2]float32{
		bbox[0] + bbox[2]/2,
		bbox[1] + bbox[3],
	}
}

// Top returns center top point
func (bbox BBox) Top() [2]float32 {
	return [2]float32{
		bbox[0],
		bbox[1] + bbox[3]/2,
	}
}

func (bbox BBox) Dist(b BBox) float32 {
	p1 := bbox.Center()
	p2 := b.Center()

	x := p2[0] - p1[0]
	y := p2[1] - p1[1]

	return float32(math.Abs(math.Sqrt(float64(x*x + y*y))))
}

func (bbox BBox) IsEqual(b BBox) bool {
	return bbox[0] == b[0] && bbox[1] == b[1] && bbox[2] == b[2] && bbox[3] == b[3]
}

// RestoreSize should restore bbox size relative crop BBox
func (bbox BBox) RestoreSize(crop BBox) BBox {
	return BBox{
		bbox[0]*crop[2] + crop[0],
		bbox[1]*crop[3] + crop[1],
		bbox[2] * crop[2],
		bbox[3] * crop[3],
	}
}
