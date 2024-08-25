package fxtools

import (
	"fmt"
	"image/color"
	"math"
	"time"
)

type HDRColor struct {
	R float64
	G float64
	B float64
	A float64
}

// How to achieve lighting effects:
// 1. Use a HDR color model (colors with values greater than 1.0)
// 2. Add up all the light sources in the scene for each cell
// 3. Apply the combined light to the cell's color by multiplying the light color with the cell's color
// 4. Use a tone mapping algorithm to convert HDR colors to LDR colors

func (h HDRColor) EncodeAsString() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", h.R, h.G, h.B)
}
func NewColorFromString(colorString string) HDRColor {
	var r, g, b float64
	fmt.Sscanf(colorString, "(%f, %f, %f)", &r, &g, &b)
	return HDRColor{r, g, b, 1.0}
}
func (h HDRColor) Multiply(color HDRColor) HDRColor {
	return HDRColor{
		R: h.R * color.R,
		G: h.G * color.G,
		B: h.B * color.B,
		A: h.A * color.A,
	}
}

func (h HDRColor) Lerp(otherColor HDRColor, percent float64) HDRColor {
	return HDRColor{
		R: Lerp(h.R, otherColor.R, percent),
		G: Lerp(h.G, otherColor.G, percent),
		B: Lerp(h.B, otherColor.B, percent),
		A: Lerp(h.A, otherColor.A, percent),
	}
}
func NewColorFromRGBA(rgba color.RGBA) HDRColor {
	return HDRColor{
		R: float64(rgba.R) / 0xFF,
		G: float64(rgba.G) / 0xFF,
		B: float64(rgba.B) / 0xFF,
		A: float64(rgba.A) / 0xFF,
	}
}
func Lerp(first float64, second float64, percent float64) float64 {
	return first + percent*(second-first)
}

func (h HDRColor) ToRGB() HDRColor {
	return h
}

func (h HDRColor) AValue() float64 {
	return h.A
}
func (h HDRColor) MultiplyWithScalar(factor float64) HDRColor {
	return HDRColor{h.R * factor, h.G * factor, h.B * factor, h.A}
}

func (h HDRColor) RValue() float64 {
	return h.R
}

func (h HDRColor) GValue() float64 {
	return h.G
}

func (h HDRColor) BValue() float64 {
	return h.B
}

func (h HDRColor) RGBA() (r, g, b, a uint32) {
	return h.ExposureToneMapping()
}

func (h HDRColor) ToRGBA() color.RGBA {
	r, g, b, a := h.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
func (h HDRColor) Brightness() float64 {
	return 0.2126*h.R + 0.7152*h.G + 0.0722*h.B
}
func (h HDRColor) ExposureToneMapping() (r, g, b, a uint32) {
	exposure := 1.0
	//lightness := R.Lightness()
	scale := float64(0xffff) // * math.Sqrt(lightness)
	// vec3 mapped = hdrColor / (hdrColor + vec3(1.0));
	//gamma := 2.2
	mappedR := 1.0 - math.Exp(-(h.R * exposure))
	r = uint32(mappedR * scale)
	mappedG := 1.0 - math.Exp(-(h.G * exposure))
	g = uint32(mappedG * scale)
	mappedB := 1.0 - math.Exp(-(h.B * exposure))
	b = uint32(mappedB * scale)
	a = uint32(h.A * scale)
	return r, g, b, a
}

func ACESFilm(x float64) float64 {
	a := 2.51
	b := 0.03
	c := 2.43
	d := 0.59
	e := 0.14
	return Clamp((x*(a*x+b))/(x*(c*x+d)+e), 0.0, 1.0)
}

func (h HDRColor) ACESFilmMapping() (r, g, b, a uint32) {

	//lightness := R.Lightness()
	scale := float64(0xffff) // * math.Sqrt(lightness)
	// vec3 mapped = hdrColor / (hdrColor + vec3(1.0));
	//gamma := 2.2
	mappedR := ACESFilm(h.R)
	r = uint32(mappedR * scale)
	mappedG := ACESFilm(h.G)
	g = uint32(mappedG * scale)
	mappedB := ACESFilm(h.B)
	b = uint32(mappedB * scale)
	a = uint32(h.A * scale)
	return r, g, b, a
}

// RelativeLuminance() is gamma-compressed
func (h HDRColor) RelativeLuminance() float64 {
	return 0.2126*h.R + 0.7152*h.G + 0.0722*h.B
}

// Luminance() is not gamma-compressed
func (h HDRColor) Luminance() float64 {
	return 0.2126*degamma(h.R) + 0.7152*degamma(h.G) + 0.0722*degamma(h.B)
}

func degamma(channelValue float64) float64 {
	// Send this function a decimal sRGB gamma encoded color value
	// between 0.0 and 1.0, and it returns a linearized value.
	if channelValue <= 0.04045 {
		return channelValue / 12.92
	} else {
		return math.Pow((channelValue+0.055)/1.055, 2.4)
	}
}

func (h HDRColor) WithClampTo(intensity float64) HDRColor {
	return HDRColor{
		R: Clamp(h.R, 0, intensity),
		G: Clamp(h.G, 0, intensity),
		B: Clamp(h.B, 0, intensity),
		A: h.A,
	}
}
func NewRGBColorFromBytes(r, g, b byte) HDRColor {
	rgbColor := HDRColor{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0, 1.0}
	return rgbColor
}
func RGBAFrom32Bits(r, g, b, a uint32) HDRColor {
	return HDRColor{
		R: float64(r) / float64(0xFFFF),
		G: float64(g) / float64(0xFFFF),
		B: float64(b) / float64(0xFFFF),
		A: float64(a) / float64(0xFFFF),
	}
}
func (h HDRColor) WithAlpha(f float64) HDRColor {
	return HDRColor{
		R: h.R,
		G: h.G,
		B: h.B,
		A: f,
	}
}

func (h HDRColor) AddRGBA(lightColor color.RGBA) HDRColor {
	return HDRColor{
		R: h.R + float64(lightColor.R)/255.0,
		G: h.G + float64(lightColor.G)/255.0,
		B: h.B + float64(lightColor.B)/255.0,
		A: h.A,
	}
}

func (h HDRColor) Add(lightColor HDRColor) HDRColor {
	return HDRColor{
		R: h.R + lightColor.R,
		G: h.G + lightColor.G,
		B: h.B + lightColor.B,
		A: h.A,
	}
}

func GetAmbientLightFromDayTime(timeOfDay time.Time) HDRColor {
	morning := HDRColor{R: 0.4, G: 0.368, B: 0.3466666666666667, A: 1.0}
	noon := HDRColor{R: 1.8, G: 1.7966666666666666, B: 1.7666666666666666, A: 1.0}
	evening := HDRColor{R: 0.2177777777777778, G: 0.2177777777777778, B: 0.26666666666666666, A: 1.0}
	night := HDRColor{R: 0.1111111111111111, G: 0.1111111111111111, B: 0.13333333333333333, A: 1.0}

	secondsSinceMidnight := timeOfDay.Hour()*3600 + timeOfDay.Minute()*60 + timeOfDay.Second()

	// we need to determine the percentage of the interval between the two times
	// for example, if it's 9:00, we need to know how far we are between 6:00 and 12:00
	// 9:00 is 50% of the way between 6:00 and 12:00
	// BUT: we want to be precise, so we need to know how many seconds are in the interval
	// 6:00 - 12:00 is 6 hours, so 6 * 3600 = 21600 seconds
	// that means we got these intervals:
	// 0:00 - 6:00 = 0 - 21600
	// 6:00 - 12:00 = 21600 - 43200
	// 12:00 - 18:00 = 43200 - 64800
	// 18:00 - 24:00 = 64800 - 86400

	var ambientLightColor HDRColor
	startColor := morning
	endColor := noon
	// 6:00 - 12:00
	if secondsSinceMidnight >= 21600 && secondsSinceMidnight < 43200 {
		intervalPercentage := float64(secondsSinceMidnight-21600) / 21600.0
		ambientLightColor = startColor.Lerp(endColor, intervalPercentage)
	} else if secondsSinceMidnight >= 43200 && secondsSinceMidnight < 64800 {
		startColor = noon
		endColor = evening
		intervalPercentage := float64(secondsSinceMidnight-43200) / 21600.0
		ambientLightColor = startColor.Lerp(endColor, intervalPercentage)
	} else if secondsSinceMidnight >= 64800 && secondsSinceMidnight < 86400 {
		startColor = evening
		endColor = night
		intervalPercentage := float64(secondsSinceMidnight-64800) / 21600.0
		ambientLightColor = startColor.Lerp(endColor, intervalPercentage)
	} else if secondsSinceMidnight >= 0 && secondsSinceMidnight < 21600 {
		startColor = night
		endColor = morning
		intervalPercentage := float64(secondsSinceMidnight) / 21600.0
		ambientLightColor = startColor.Lerp(endColor, intervalPercentage)
	}
	return ambientLightColor
}
