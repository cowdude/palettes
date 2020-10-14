//go:generate go run .
package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/cowdude/palettes"
)

var (
	width     = flag.Int("width", 640, "preview image width")
	height    = flag.Int("height", 480, "preview image height")
	minColors = flag.Int("min-colors", 256, "minimum number of colors per palette")
)

func init() {
	flag.Parse()
}

func density(x, y float64) float64 {
	return 0.5 + 0.5*math.Sin(13*(math.Pow(x, 3)+math.Pow(y, 2)))
}

func main() {
	w, h := *width, *height
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	f, err := os.Create("example.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	build := func(baseColors []palettes.RGBAF64) palettes.Palette {
		return palettes.Definition{
			Base:      baseColors,
			MinColors: *minColors,
		}.Build()
	}

	allPalettes := [...]palettes.Palette{
		build(palettes.Inferno),
		build(palettes.Magma),
		build(palettes.Viridis),
		build(palettes.RdGy),
	}
	const plotSpan = 2

	dx := w / plotSpan
	dy := h / plotSpan

	for i, palette := range allPalettes {
		x0 := (i % plotSpan) * dx
		y0 := (i / plotSpan) * dy

		for y := y0; y < y0+dy; y++ {
			for x := x0; x < x0+dx; x++ {
				u := float64(x-x0) / float64(dx-1)
				v := float64(y-y0) / float64(dy-1)
				z := density(u, v)
				sample := palette.Sample(z)
				img.Set(x, y, sample)
			}
		}
	}

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	fmt.Println("Generated:", f.Name())
}
