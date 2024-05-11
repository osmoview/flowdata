package flowdata

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckObjectLineCrossing(t *testing.T) {
	type input struct {
		objects []*Object
		a       [2]float64
		b       [2]float64
		anchor  AnchorType
	}

	type expected struct {
		crossedOneWay     []*Object
		crossedAnotherWay []*Object
		notCrossed        []*Object
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "ok_bottom",
			input: input{
				objects: []*Object{
					{
						ID:        "id",
						Type:      "type",
						Name:      "name",
						Prob:      1,
						BBox:      [4]float64{0.1, 0.1, 0.2, 0.2},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
				},
				a:      [2]float64{0.4, 0.4},
				b:      [2]float64{0.6, 0.6},
				anchor: Bottom,
			},
			expected: expected{
				crossedOneWay: []*Object{},
				crossedAnotherWay: []*Object{
					{
						ID:        "id",
						Type:      "type",
						Name:      "name",
						Prob:      1,
						BBox:      [4]float64{0.1, 0.1, 0.2, 0.2},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
				},
				notCrossed: []*Object{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nl := NewLine("test", test.input.a, test.input.b, color.RGBA{})

			crossedOneWay, crossedAnotherWay, notCrossed := nl.CheckObjectLineCrossing(test.input.objects, test.input.anchor)

			switch test.name {
			case "ok_bottom":
				require.Equal(t, test.expected.crossedOneWay, crossedOneWay)
				require.Equal(t, test.expected.crossedAnotherWay, crossedAnotherWay)
				require.Equal(t, test.expected.notCrossed, notCrossed)
			}
		})
	}
}
