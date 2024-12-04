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
	Text  string       `json:"text,omitempty"`
	Color string       `json:"color,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`

	Keypoints Keypoints `json:"keypoints,omitempty"`

	Objects []*Object `json:"objects,omitempty"`

	ExternalID  string    `json:"external_id,omitempty"`
	Vector      Vector    `json:"vector,omitempty"`
	VectorMatch []float32 `json:"vector_match,omitempty"`
	Profile     *Profile  `json:"profile,omitempty"`
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

	// copy data if it exists
	if o.Data != nil {
		newObj.Data = make(map[string]interface{})
		for k, v := range o.Data {
			newObj.Data[k] = v
		}
	}

	// create copies of sub objects
	if o.Objects != nil {
		newObj.Objects = make([]*Object, 0, len(o.Objects))
		for _, sub := range o.Objects {
			newObj.Objects = append(newObj.Objects, sub.Copy())
		}
	}

	return &newObj
}

// Keypoints is a slice on keypoints
type Keypoints []Keypoint

// Keypoint is a point with name and may has links to another keypoints
type Keypoint struct {
	Key   string     `json:"key"`
	Point [2]float64 `json:"point"`
	Links []string   `json:"links"`
}
