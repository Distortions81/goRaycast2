package main

import (
	"image/color"

	"github.com/chewxy/math32"
)

// Linearize sRGB to linear space (remove gamma correction)
func sRGBToLinear(value float32) float32 {
	if value <= 0.04045 {
		return value / 12.92
	}
	return math32.Pow((value+0.055)/1.055, 2.4)
}

// Apply gamma correction to convert from linear space to sRGB
func linearTosRGB(value float32) float32 {
	if value <= 0.0031308 {
		return 12.92 * value
	}
	return 1.055*math32.Pow(value, 1.0/2.4) - 0.055
}

// Calculate light falloff using inverse-square law
func lightFalloff(distance float32, intensity float32) float32 {
	// Basic inverse-square law: falloff = intensity / (distance^2)
	if distance <= 0 {
		return intensity // Prevent division by zero
	}
	return intensity / (distance * distance)
}

// Calculate color with falloff and gamma correction
func applyFalloff(distance float32, intensity float32, value float32) float32 {

	// Linearize each color channel
	linear := sRGBToLinear(value)

	// Apply light falloff
	falloff := lightFalloff(distance, intensity)

	// Multiply each linear channel by falloff
	linear *= falloff

	return (math32.Max(0, math32.Min(1, linearTosRGB(linear))))
}

// HSVtoRGB converts HSV values to RGB
func HSVtoRGB(h, s, v float32) color.NRGBA {
	c := v * s
	x := c * (1 - math32.Abs(math32.Mod(h/60.0, 2)-1))
	m := v - c

	var r1, g1, b1 float32

	if h >= 0 && h < 60 {
		r1, g1, b1 = c, x, 0
	} else if h >= 60 && h < 120 {
		r1, g1, b1 = x, c, 0
	} else if h >= 120 && h < 180 {
		r1, g1, b1 = 0, c, x
	} else if h >= 180 && h < 240 {
		r1, g1, b1 = 0, x, c
	} else if h >= 240 && h < 300 {
		r1, g1, b1 = x, 0, c
	} else if h >= 300 && h < 360 {
		r1, g1, b1 = c, 0, x
	}

	// Convert to RGB by adding m and scaling to the range of 0-255
	r := (r1 + m) * 255
	g := (g1 + m) * 255
	b := (b1 + m) * 255

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 1}
}
