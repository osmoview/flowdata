package flowdata

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

type TransformParams struct {
	Rotate int
	Flip   string // "", "horizontal", "vertical", "both"
}

func (p TransformParams) NeedsTransform() bool {
	return p.Rotate != 0 || p.Flip != ""
}

// ApplyTransform applies rotation and flip transformations to the frame
func applyTransform(f *Frame, params TransformParams) (*Frame, error) {
	var err error

	if params.Rotate != 0 {
		f, err = f.Rotate(params.Rotate)
		if err != nil {
			return nil, err
		}
	}

	if params.Flip != "" {
		f, err = f.Flip(params.Flip)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

// Transform applies rotation and flip transformations to the frame
func (f *Frame) Transform(params TransformParams) (*Frame, error) {
	if !params.NeedsTransform() {
		return f, nil
	}

	return applyTransform(f, params)
}

// Rotate applies rotation transformation to the frame
func (f *Frame) Rotate(angle int) (*Frame, error) {
	if angle == 0 {
		return f, nil
	}

	rotationAngle := float64(-angle)
	img := imaging.Rotate(f.GetImage(), rotationAngle, image.Black)
	return NewFrameImage(img)
}

// Flip applies flip transformation to the frame
// direction: "horizontal", "vertical", or "both"
func (f *Frame) Flip(direction string) (*Frame, error) {
	if direction == "" {
		return f, nil
	}

	var img image.Image

	switch direction {
	case "horizontal":
		img = imaging.FlipH(f)
	case "vertical":
		img = imaging.FlipV(f)
	case "both":
		img = imaging.FlipH(imaging.FlipV(f))
	default:
		return nil, fmt.Errorf("invalid flip direction: %s (valid: horizontal, vertical, both)", direction)
	}

	return NewFrameImage(img)
}
