package captchy

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"

	"golang.org/x/image/font/gofont/gobold"

	"golang.org/x/image/font/gofont/gomono"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	rulerDeviation            = 20
	baselineDeviation         = 8
	rotateDeviation   float64 = 5
	sizeDeviation     float64 = 5
	amplude           float64 = 7
	period            float64 = 150
	maxGridApply              = 5
	circleCount               = 1
	maxNoiseLine              = 5
	saltPercent               = 100 >> 3
)

var (
	imgW        int
	imgH        int
	capLen      int
	tFont       *truetype.Font
	size        float64 //font size in points
	fontColors  []color.Color
	RulerColor  = ColorInHex("0x456D7C")
	GridColor   = ColorInHex("0xB9B6B6")
	CircleColor = ColorInHex("#728C7F")
	noiseColors []color.Color
	bgSrc       *image.Uniform
)

const (
	DefaultSize float64 = 48
	DefaultDPI  float64 = 72

	RegularFont = string(iota)
	MonoFont
	BoldFont
	MediumFont
	ItalicFont
)

type Config struct {
	ImgW        int
	ImgH        int
	Len         int
	fontPath    string
	FontSize    float64
	FontColors  []color.Color
	BgColor     color.Color
	NoiseColors []color.Color
}

func ColorInHex(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "0x")
	hex = strings.TrimPrefix(hex, "#")

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		log.Fatalln(err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		log.Fatalln(err)
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF}
}

func Default() *Config {
	cfg := Config{
		240,
		80,
		6,
		RegularFont,
		40,
		[]color.Color{color.Black, ColorInHex("0xFA9D9E"), ColorInHex("0xF8AA5D"), ColorInHex("0x429BF5"), ColorInHex("0xFA6565")},
		color.White,
		[]color.Color{
			ColorInHex("0x3FBFBF"),
			ColorInHex("0x121DF1"),
			ColorInHex("0xF1AA12"),
		},
	}

	return &cfg
}

func New(cfg *Config) {
	switch cfg.fontPath {
	case RegularFont:
		tFont, _ = truetype.Parse(goregular.TTF)
	case MonoFont:
		tFont, _ = truetype.Parse(gomono.TTF)
	case BoldFont:
		tFont, _ = truetype.Parse(gobold.TTF)
	case MediumFont:
		tFont, _ = truetype.Parse(gomedium.TTF)
	case ItalicFont:
		tFont, _ = truetype.Parse(goitalic.TTF)
	default:
		fontBytes, err := ioutil.ReadFile(cfg.fontPath)
		if err != nil {
			log.Fatalln(err)
		}
		tFont, err = truetype.Parse(fontBytes)
		if err != nil {
			log.Fatalln(err)
		}
	}

	imgW = cfg.ImgW
	imgH = cfg.ImgH
	capLen = cfg.Len
	size = cfg.FontSize
	bgSrc = image.NewUniform(cfg.BgColor)
	noiseColors = cfg.NoiseColors
	fontColors = cfg.FontColors
}

func (c *Config) SetFont(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalln("Invalid font file path")
	}
	c.fontPath = path
}

func SetNoiseThreshold(p float64) {
	if p > 0 && p < 1 {
		np := int(p * 100)
		saltPercent = np
	} else if p == 0 {
		saltPercent = 0
		maxNoiseLine = 0
		maxGridApply = 0
	}
}

func SetBg(c color.Color) {
	bgSrc = image.NewUniform(c)
}

func SetRulerDeviation(d int) {
	if d > (imgH >> 1) {
		log.Fatalln("Ruler deviation should not greater than Height / 2")
	}
	rulerDeviation = d
}

func SetFontDeviation(d int) {
	if d > (imgH >> 1) {
		log.Fatalln("Font deviation should not greater than Height / 2")
	}
	baselineDeviation = d
}

func SetRotateDegree(d float64) {
	if d >= 0 && d < 8 {
		rotateDeviation = d
		return
	}
	log.Fatalln("Rotate degree should be in range [0, 8)")
}

func DisableRuler() {
	RulerColor = color.Transparent
}

func DisableCircle() {
	circleCount = 0
}

func GenerateImg(t []byte) Encoder {
	d := newDrawer()
	if saltPercent > 0 && len(noiseColors) == 0 {
		d.ApplySaltEffect()
	}

	for i := 0; i < circleCount; i++ {
		d.DrawCircle(rand.Intn(imgW), rand.Intn(imgH), 34+devateRandI(10), CircleColor)
	}

	for i := 0; i < maxGridApply; i++ {
		d.ApplyGridEffect()
	}
	d.DrawNoiseLines()

	if RulerColor != color.Transparent {
		d.DrawRuler()
	}
	d.DrawText(t)
	d.ApplyDistort(amplude, period)

	if rotateDeviation > 0 {
		d.ApplyRotation()
	}

	return &CaptchaEncoder{d.Image()}
}
