package palettes

import (
	"image/color"
	"math"
)

type RGBAF64 struct {
	R, G, B, A float64
}

type Palette []color.RGBA
type BaseColors []RGBAF64

type Definition struct {
	Base      BaseColors
	MinColors int
}

func (f64 RGBAF64) RGBA() color.RGBA {
	conv := func(x float64) uint8 {
		return uint8(math.Max(0, math.Min(0xFF, x*0xFF)))
	}
	return color.RGBA{
		R: conv(f64.R),
		G: conv(f64.G),
		B: conv(f64.B),
		A: conv(f64.A),
	}
}

func lerpColor(from, to RGBAF64, k float64) RGBAF64 {
	lerp := func(from, to float64) float64 {
		return float64(from) + (float64(to)-float64(from))*k
	}
	return RGBAF64{
		R: lerp(from.R, to.R),
		G: lerp(from.G, to.G),
		B: lerp(from.B, to.B),
		A: lerp(from.A, to.A),
	}
}

func rescale(data []RGBAF64, s float64) []RGBAF64 {
	for i := range data {
		data[i].R *= s
		data[i].G *= s
		data[i].B *= s
		data[i].A *= s
	}
	return data
}

func (def Definition) Build() (out Palette) {
	n := def.MinColors
	data := def.Base
	if n < len(data) {
		n = len(data)
	}
	out = make(Palette, n)

	for i := range out {
		progress := float64(i) * float64(len(data)-1) / float64(n-1) // [0; 1]
		ipart, fpart := math.Modf(progress)
		j := int(ipart)
		if j == len(data)-1 {
			out[i] = data[j].RGBA()
		} else {
			out[i] = lerpColor(data[j], data[j+1], fpart).RGBA()

		}
	}
	return
}

func (palette Palette) Sample(normVal float64) color.RGBA {
	n := int(normVal * float64(len(palette)-1))
	if lim := len(palette) - 1; n > lim {
		n = lim
	}
	if n < 0 {
		n = 0
	}
	return palette[n]
}
