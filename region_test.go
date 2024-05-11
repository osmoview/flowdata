package flowdata

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckEntranceInRegion(t *testing.T) {
	type input struct {
		objects []*Object
		points  [4]float64
	}

	type expected struct {
		in  []*Object
		out []*Object
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "ok",
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
					{
						ID:        "id2",
						Type:      "type2",
						Name:      "name2",
						Prob:      1,
						BBox:      [4]float64{0.2, 0.2, 0.1, 0.1},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
					{
						ID:        "id3",
						Type:      "type3",
						Name:      "name3",
						Prob:      1,
						BBox:      [4]float64{0.6, 0.6, 0.3, 0.3},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
					{
						ID:        "id4",
						Type:      "type4",
						Name:      "name4",
						Prob:      1,
						BBox:      [4]float64{0.6, 0.6, 0.2, 0.2},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
				},
				points: [4]float64{0.0, 0.5, 0.0, 0.5},
			},
			expected: expected{
				in: []*Object{
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
					{
						ID:        "id2",
						Type:      "type2",
						Name:      "name2",
						Prob:      1,
						BBox:      [4]float64{0.2, 0.2, 0.1, 0.1},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
				},
				out: []*Object{
					{
						ID:        "id3",
						Type:      "type3",
						Name:      "name3",
						Prob:      1,
						BBox:      [4]float64{0.6, 0.6, 0.3, 0.3},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
					{
						ID:        "id4",
						Type:      "type4",
						Name:      "name4",
						Prob:      1,
						BBox:      [4]float64{0.6, 0.6, 0.2, 0.2},
						Area:      [][2]float64{},
						Data:      map[string]interface{}{},
						Keypoints: Keypoints{},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nr := NewRegion("test", test.input.points, color.RGBA{})

			in, out := nr.CheckEntranceInRegion(test.input.objects)

			switch test.name {
			case "ok":
				require.Equal(t, test.expected.in, in)
				require.Equal(t, test.expected.out, out)
			}
		})
	}
}

func TestBBoxInRegion(t *testing.T) {
	type input struct {
		bbox   [4]float64
		points [4]float64
	}

	tests := []struct {
		name     string
		input    input
		expected bool
	}{
		{
			name: "ok_box_in_region",
			input: input{
				bbox:   [4]float64{0.1, 0.1, 0.2, 0.2},
				points: [4]float64{0.0, 0.5, 0.0, 0.5},
			},
			expected: true,
		},
		{
			name: "ok_box_not_in_region",
			input: input{
				bbox:   [4]float64{0.6, 0.6, 0.3, 0.3},
				points: [4]float64{0.0, 0.5, 0.0, 0.5},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nr := NewRegion("test", test.input.points, color.RGBA{})

			result := nr.bBoxInRegion(test.input.bbox, nr.Points)

			switch test.name {
			case "ok_box_in_region":
				require.Equal(t, test.expected, result)
			case "ok_box_not_in_region":
				require.Equal(t, test.expected, result)
			}
		})
	}
}

func TestCalculateBBoxCenter(t *testing.T) {
	type expected struct {
		centerX float64
		centerY float64
	}

	tests := []struct {
		name     string
		input    [4]float64
		expected expected
	}{
		{
			name:  "ok",
			input: [4]float64{0.1, 0.1, 0.2, 0.2},
			expected: expected{
				centerX: 0.2,
				centerY: 0.2,
			},
		},
		{
			name:  "err",
			input: [4]float64{0.1, 0.1, 0.2, 0.2},
			expected: expected{
				centerX: 0.3,
				centerY: 0.3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			centerX, centerY := CalculateBBoxCenter(test.input)

			switch test.name {
			case "ok":
				require.Equal(t, test.expected.centerX, centerX)
				require.Equal(t, test.expected.centerY, centerY)
			case "err":
				require.NotEqual(t, test.expected.centerX, centerX)
				require.NotEqual(t, test.expected.centerY, centerY)
			}
		})
	}
}
