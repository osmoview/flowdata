package flowdata

import "fmt"

type Object struct {
	ID   string  `json:"id,omitmepty"`
	Type string  `json:"type"`
	Name string  `json:"name,omitmepty"`
	Prob float64 `json:"prob,omitmepty"`
	// rectangle format [x, y, w, h] displays where an object is located
	BBox [4]float64 `json:"bbox,omitmepty"`
	// format [ [x, y], [x, y] ... ] to highlight the area occupied by an object in the frame
	Area [][2]float64 `json:"area,omitmepty"`

	Data map[string]interface{} `json:"data,omitmepty"`

	// for example a person, his keypoints are his head, shoulders, arms, etc.
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
