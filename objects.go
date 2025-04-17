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
	Items   []*Object `json:"items,omitempty"` // DEPRECATED: required only to compatible with legacy services

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

// Merge two objects
func (o *Object) Merge(obj *Object) {
	if o.Type != "" && o.Type != obj.Type {
		return // can not merge different types
	}

	if obj.ID != "" {
		o.ID = obj.ID
	}
	if obj.Name != "" {
		o.Name = obj.Name
	}
	if obj.Prob != 0 {
		o.Prob = obj.Prob
	}
	if obj.BBox.IsValid() {
		o.BBox = obj.BBox
	}
	if len(obj.Area) > 0 {
		o.Area = obj.Area
	}
	if obj.Text != "" {
		o.Text = obj.Text
	}
	if obj.Color != "" {
		o.Color = obj.Color
	}

	if obj.Data != nil {
		if o.Data == nil {
			o.Data = obj.Data
		} else {
			for k, v := range obj.Data {
				o.Data[k] = v
			}
		}
	}

	if len(obj.Keypoints) != 0 {
		o.Keypoints = obj.Keypoints
	}

	o.Objects = MergeObjects(o.Objects, obj.Objects)
	o.Items = MergeObjects(o.Items, obj.Objects)

	if obj.ExternalID != "" {
		o.ExternalID = obj.ExternalID
	}

	if len(obj.Vector) != 0 {
		o.Vector = obj.Vector
		o.VectorMatch = obj.VectorMatch
	}

	if obj.Profile != nil {
		o.Profile = obj.Profile
	}

}

// Merge two slices of objects
func MergeObjects(a, b []*Object) []*Object {
	if len(b) == 0 && len(a) == 0 {
		return a
	}
	if len(b) == 0 && len(a) > 0 {
		return a
	}
	if len(b) > 0 && len(a) == 0 {
		return b
	}

	list := make([]*Object, 0, len(a))
	for _, o := range a {
		obj, ok := LookupSimilarObject(b, o)
		if ok {
			o.Merge(obj)
		}
		list = append(list, o)
	}

	for _, o := range b {
		if _, ok := LookupSimilarObject(list, o); !ok {
			list = append(list, o)
		}
	}

	return list
}

func LookupSimilarObject(list []*Object, o *Object) (*Object, bool) {
	for _, obj := range list {
		if obj.ID != "" && obj.ID == o.ID {
			return obj, true
		}

		if obj.BBox.IsValid() && obj.BBox.IsEqual(o.BBox) {
			return obj, true
		}
	}

	return nil, false
}

// Keypoints is a slice on keypoints
type Keypoints []Keypoint

// Keypoint is a point with name and may has links to another keypoints
type Keypoint struct {
	Key   string     `json:"key"`
	Point [2]float64 `json:"point"`
	Links []string   `json:"links"`
}
