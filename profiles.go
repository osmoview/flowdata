package flowdata

import (
	"log"
	"time"
)

type Profile struct {
	// uuid
	ID string `json:"id"`

	CreatedAt time.Time `json:"created_at"`

	// universal profile name, like person firstname + surname etc, or name of car, etc
	Name string `json:"name"`

	// profiles type id
	Type string `json:"type"`

	// some external ID, like car number, some code etc
	ExternalID string `json:"external_id,omitempty"`

	// some territory to which this profile is related
	Territory string `json:"territory,omitempty"`

	// custom data
	Data map[string]any `json:"data,omitempty" gorm:"serializer:json;type:json"`

	Images []ProfileImage `json:"images,omitempty" gorm:"serializer:json;type:json"`
}

func (p *Profile) TableName() string {
	return "profiles"
}

type ProfileImage struct {
	Filename string `json:"filename"`
	Vector   Vector `json:"vector,omitempty"`
}

func (p Profile) VectorMatch(v Vector, max, min float32) (dist float32, iter int, ok bool) {
	if min > max {
		min, max = max, min
	}

	for i, pimg := range p.Images {
		// if profiles has image but without vector, ignore it
		if !pimg.Vector.IsValid() {
			continue
		}

		// calculate distance between 2 vectors
		d, err := pimg.Vector.EuclideanDistance(v)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		// if distance less max treshold, this means that object can match this profile
		if d < max {
			// if distance less min threshold, so this object is this profile
			if d < min {
				return d, i, true
			}

			if dist == 0 || d < dist {
				dist = d
				iter = i
			}
		}
	}

	return dist, iter, dist > 0
}
