package captchy

import (
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

func Test_GenPNGfile(t *testing.T) {
	New(Default())
	rbs := RandomString()
	encoder := GenerateImg(rbs)

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	encoder.ToPNG(outFile)
}

func Test_GenJPGfile(t *testing.T) {
	New(Default())
	rbs := RandomString()
	encoder := GenerateImg(rbs)

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.jpeg")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	ops := &jpeg.Options{Quality: jpeg.DefaultQuality}
	encoder.ToJPEG(outFile, ops)
}

func Test_GenSimple(t *testing.T) {
	cfg := &Config{
		240,
		80,
		4,
		ItalicFont,
		40,
		[]color.Color{color.Black},
		color.White,
		[]color.Color{},
	}
	cfg.SetFont("./testdata/luxisr.ttf")
	New(cfg)
	SetNoiseThreshold(0)
	DisableRuler()
	DisableCircle()
	SetRotateDegree(0)

	rbs := RandomString()
	encoder := GenerateImg(rbs)

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	encoder.ToPNG(outFile)
}

func BenchmarkImgGenLen6Colorful(b *testing.B) {
	New(Default())
	rbs := RandomString()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateImg(rbs)
	}
}

func BenchmarkImgGenLen4Normal(b *testing.B) {
	New(&Config{
		240,
		80,
		4,
		RegularFont,
		40,
		[]color.Color{color.Black},
		color.White,
		[]color.Color{
			ColorInHex("0x3FBFBF"),
			ColorInHex("0x121DF1"),
			ColorInHex("0xF1AA12"),
		},
	})
	rbs := RandomString()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateImg(rbs)
	}
}
