package flowdata

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

func (bbox BBox) Center() [2]float64 {
	return [2]float64{
		bbox[0] + bbox[2]/2,
		bbox[1] + bbox[3]/2,
	}
}

func (bbox BBox) Bottom() [2]float64 {
	return [2]float64{
		bbox[0] + bbox[2]/2,
		bbox[1] + bbox[3],
	}
}

func (bbox BBox) Top() [2]float64 {
	return [2]float64{
		bbox[0],
		bbox[1] + bbox[3]/2,
	}
}
