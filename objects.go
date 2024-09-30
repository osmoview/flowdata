package flowdata

import (
	"fmt"
)

type Object struct {
	ID    string       `json:"id,omitempty"`
	Type  string       `json:"type"`
	Name  string       `json:"name,omitempty"`
	Prob  float64      `json:"prob,omitempty"`
	BBox  BBox         `json:"bbox,omitempty"`
	Area  [][2]float64 `json:"area,omitempty"`
	Color string       `json:"color,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`

	Keypoints Keypoints `json:"keypoints,omitempty"`
}

func (o *Object) String() string {
	if o.Name != "" {
		return fmt.Sprintf("%s:%s", o.Type, o.Name)
	}

	return o.Type
}

// IsValid check is object is valid, `type` is required field
func (o *Object) IsValid() bool {
	return o.Type != ""
}

// Copy creates and returns a copy of the Object
func (o *Object) Copy() *Object {
	newObj := *o

	newObj.Area = make([][2]float64, len(o.Area))
	copy(newObj.Area, o.Area)

	if o.Data != nil {
		newObj.Data = make(map[string]interface{})
		for k, v := range o.Data {
			newObj.Data[k] = v
		}
	}

	return &newObj
}

// IsZeroBBox check is empty bbox
func (o *Object) IsZeroBBox() bool {
	for _, v := range o.BBox {
		if v != 0 {
			return false
		}
	}
	return true
}

// Keypoints is a slice on keypoints
type Keypoints []Keypoint

// Keypoint is a point with name and may has links to another keypoints
type Keypoint struct {
	Key   string     `json:"key"`
	Point [2]float64 `json:"point"`
	Links []string   `json:"links"`
}

type BBox [4]float64

func (bbox BBox) Scale(n float64) BBox {
	kw := bbox[2]*n - bbox[2]
	kh := bbox[3]*n - bbox[3]

	return BBox{
		bbox[0] - kw/2,
		bbox[1] - kh/2,
		bbox[2] + kw,
		bbox[3] + kh,
	}
}
