package flowdata

import "fmt"

type Object struct {
	ID    string       `json:"id,omitempty"`
	Type  string       `json:"type"`
	Name  string       `json:"name,omitempty"`
	Prob  float64      `json:"prob,omitempty"`
	BBox  [4]float64   `json:"bbox,omitempty"`
	Area  [][2]float64 `json:"area,omitempty"`
	Color string       `json:"color,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`

	Keypoints Keypoints `json:"keypoints,omitempty"`
}

func (o Object) String() string {
	if o.Name != "" {
		return fmt.Sprintf("%s:%s", o.Type, o.Name)
	}

	return o.Type
}

// IsValid check is object is valid, `type` is required field
func (o Object) IsValid() bool {
	return o.Type != ""
}

// Keypoints is a slice on keypoints
type Keypoints []Keypoint

// Keypoint is a point with name and may has links to another keypoints
type Keypoint struct {
	Key   string
	Point [2]float64
	Links []string
}
