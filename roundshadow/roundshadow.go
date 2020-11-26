package roundshadow

import (
	"image"
	"image/color"
	"math"
)

func NewRoundShadow(img image.Image, p image.Point, d int) RoundShadow {
	return RoundShadow{img, p, d}
}

type RoundShadow struct {
	image    image.Image
	point    image.Point
	diameter int
}

func (r RoundShadow) ColorModel() color.Model {
	return r.image.ColorModel()
}

func (r RoundShadow) Bounds() image.Rectangle {
	return image.Rect(0, 0, r.diameter, r.diameter)
}

func (r RoundShadow) At(x, y int) color.Color {
	d := r.diameter
	dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
	if dis > float64(d)/2 {
		return r.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
	} else {
		return r.image.At(r.point.X+x, r.point.Y+y)
	}
}
