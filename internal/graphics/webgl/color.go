//go:build js

package webgl

type Color struct {
	R, G, B, A float64
}

func ColorBlack(a float64) *Color {
	return &Color{
		R: 0,
		G: 0,
		B: 0,
		A: a,
	}
}

func ColorWhite(a float64) *Color {
	return &Color{
		R: 1,
		G: 1,
		B: 1,
		A: a,
	}
}

func ColorRed(a float64) *Color {
	return &Color{
		R: 1,
		G: 0,
		B: 0,
		A: a,
	}
}

func ColorGreen(a float64) *Color {
	return &Color{
		R: 0,
		G: 1,
		B: 0,
		A: a,
	}
}

func ColorBlue(a float64) *Color {
	return &Color{
		R: 0,
		G: 0,
		B: 1,
		A: a,
	}
}
