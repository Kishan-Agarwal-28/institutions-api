package main

import (
	"crypto/md5"
	"fmt"
	"math"
)

// clamp restricts a value between 0 and 1
func clamp(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// toByte converts a float value to a uint8 byte (0-255)
func toByte(f float64) uint8 {
	val := f
	if val <= 0.0031308 {
		val = 12.92 * val
	} else {
		val = 1.055*math.Pow(val, 1.0/2.4) - 0.055
	}
	return uint8(clamp(val) * 255)
}

// oklchToHex converts OKLCH color space to hex color
func oklchToHex(l, c, h float64) string {
	hRad := h * (math.Pi / 180.0)
	a := c * math.Cos(hRad)
	b := c * math.Sin(hRad)

	l_ := l + 0.3963377774*a + 0.2158037573*b
	m_ := l - 0.1055613458*a - 0.0638541728*b
	s_ := l - 0.0894841775*a - 1.2914855480*b

	l_ = math.Pow(l_, 3)
	m_ = math.Pow(m_, 3)
	s_ = math.Pow(s_, 3)

	red := 4.0767416621*l_ - 3.3077115913*m_ + 0.2309699292*s_
	green := -1.2684380046*l_ + 2.6097574011*m_ - 0.3413193965*s_
	blue := -0.0041960863*l_ - 0.7034186147*m_ + 1.7076147010*s_

	return fmt.Sprintf("#%02x%02x%02x", toByte(red), toByte(green), toByte(blue))
}

// generateColor creates a deterministic color from a string input
func generateColor(s string) (string, [16]byte) {
	hash := md5.Sum([]byte(s))
	hue := float64(uint16(hash[0])<<8|uint16(hash[1])) / 65535.0 * 360.0

	chroma := 0.10 + (float64(hash[2])/255.0)*0.06

	lightness := 0.55 + (float64(hash[3])/255.0)*0.13

	hex := oklchToHex(lightness, chroma, hue)
	return hex, hash
}
