package captchy

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type drawer struct {
	captchaImg *image.RGBA
	fontDrawer *font.Drawer
}

func newDrawer() *drawer {
	d := new(drawer)
	rgba := canvasPool.Get().(*image.RGBA)
	draw.Draw(rgba, rgba.Bounds(), bgSrc, image.ZP, draw.Src)
	d.captchaImg = rgba
	d.PrepareDrawer()

	return d
}

func (d *drawer) PrepareDrawer() {
	d.fontDrawer = &font.Drawer{
		Dst: d.captchaImg,
		Src: image.Black,
		Face: truetype.NewFace(tFont, &truetype.Options{
			Size:    size,
			DPI:     DefaultDPI,
			Hinting: font.HintingNone,
		}),
	}
}

func (d *drawer) DrawText(t []byte) {
	y := 10 + int(math.Ceil(size*DefaultDPI/72))
	dx := (fixed.I(imgW) - d.fontDrawer.MeasureBytes(t)).Floor() / 2
	d.fontDrawer.Dot = fixed.P(dx, y)
	gd := (d.fontDrawer.MeasureBytes(t) / fixed.I(len(t))).Floor()

	unf := image.NewUniform(fontColors[0])
	if len(fontColors) == 1 {
		d.fontDrawer.Src = unf
		for i := 0; i < len(t); i++ {
			d.fontDrawer.Dot.Y = fixed.I(y) + fixed.I(devateRandI(baselineDeviation))
			d.fontDrawer.DrawBytes(t[i : i+1])
			d.fontDrawer.Dot.X += fixed.I(gd)
		}
	} else {
		fmask := make([]int, capLen)
		for i := 0; i < len(fmask); i++ {
			fmask[i] = rand.Intn(len(fontColors))
		}

		for i := 0; i < len(t); i++ {
			unf.C = fontColors[fmask[i]]
			d.fontDrawer.Src = unf
			d.fontDrawer.Dot.Y = fixed.I(y) + fixed.I(devateRandI(baselineDeviation))
			d.fontDrawer.DrawBytes(t[i : i+1])
			d.fontDrawer.Dot.X += fixed.I(gd)
		}
	}
}

func (d *drawer) DrawRuler() {
	my := imgH>>1 + devateRandI(rulerDeviation)
	for i := 5; i < imgW-5; i++ {
		d.captchaImg.Set(i, my, RulerColor)
		d.captchaImg.Set(i+1, my, RulerColor)
	}
}

func (d *drawer) DrawCircle(x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r << 1)

	for x > y {
		d.captchaImg.Set(x0+x, y0+y, c)
		d.captchaImg.Set(x0+y, y0+x, c)
		d.captchaImg.Set(x0-y, y0+x, c)
		d.captchaImg.Set(x0-x, y0+y, c)
		d.captchaImg.Set(x0-x, y0-y, c)
		d.captchaImg.Set(x0-y, y0-x, c)
		d.captchaImg.Set(x0+y, y0-x, c)
		d.captchaImg.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r << 1)
		}
	}
}

func (d *drawer) DrawNoiseLines() {
	var x, y, x1, y1 int
	for i := 0; i < maxNoiseLine; i++ {
		x = rand.Intn(imgW)
		x1 = rand.Intn(imgW)
		y = rand.Intn(imgH)
		y1 = rand.Intn(imgH)

		drawLine(d.captchaImg, x, y, x1, y1, GridColor)
	}
}

func (d *drawer) ApplyDistort(amplude, period float64) {
	nrgba := canvasPool.Get().(*image.RGBA)
	draw.Draw(nrgba, nrgba.Bounds(), bgSrc, image.ZP, draw.Src)

	trans := color.RGBA{}

	dx := 2.0 * math.Pi / period
	for x := 0; x < imgW; x++ {
		for y := 0; y < imgH; y++ {
			if d.captchaImg.At(x, y) == bgSrc {
				continue
			}
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			rgb := d.captchaImg.At(x+int(xo), y+int(yo))
			if rgb == trans {
				nrgba.Set(x, y, image.White)
				continue
			}
			nrgba.Set(x, y, rgb)
		}
	}

	canvasPool.Put(d.captchaImg)
	d.captchaImg = nrgba
}

func (d *drawer) ApplyGridEffect() {
	gx := rand.Intn(imgW)
	gy := rand.Intn(imgH)
	gridSize := rand.Intn(20) + 10

	// compressDeviation parameter
	cd := rand.Intn(3) + 3
	for dy := 0; dy < gridSize; dy += cd {
		drawLine(d.captchaImg, gx, gy+dy, gx+gridSize, gy+dy, GridColor)
	}

	for dx := 0; dx < gridSize; dx += cd {
		drawLine(d.captchaImg, gx+dx, gy, gx+dx, gy+gridSize, GridColor)
	}
}

func (d *drawer) ApplySaltEffect() {
	cs := []int{rand.Intn(len(noiseColors)), rand.Intn(len(noiseColors))}

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			noiseCheck := rand.Intn(100)
			if noiseCheck < saltPercent {
				d.captchaImg.Set(x, y, noiseColors[x%len(cs)])
			}
		}
	}
}

func (d *drawer) ApplyRotation() {
	rotate(d.captchaImg, randAngle(rotateDeviation))
}

func (d *drawer) Image() image.Image {
	return d.captchaImg
}

func applyDistort(m, newm *image.RGBA, amplude float64, period float64) {
	trans := color.RGBA{}

	dx := 2.0 * math.Pi / period
	for x := 0; x < imgW; x++ {
		for y := 0; y < imgH; y++ {
			if m.At(x, y) == image.White {
				continue
			}
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			rgb := m.At(x+int(xo), y+int(yo))
			if rgb == trans {
				newm.Set(x, y, image.White)
				continue
			}
			newm.Set(x, y, rgb)
		}
	}
}

func rotate(img *image.RGBA, angle float64) {
	k := 4
	l := 2

	w := img.Bounds().Dx() / k
	h := img.Bounds().Dy() / l

	//fmt.Println(w, h)
	var dx int
	dy := -h

	nrgba := image.NewRGBA(image.Rect(0, 0, w, h))

	for i := 0; i < k*l; i++ {
		dx = w * (i % k)
		if i%k == 0 {
			dy += h
		}

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				hwidth := w >> 1
				hheight := h >> 1

				xt := x - hwidth
				yt := y - hheight

				sinma := math.Sin(-angle)
				cosma := math.Sqrt(1 - sinma*sinma)

				xs := int(math.Round((cosma*float64(xt)-sinma*float64(yt))+float64(hwidth))) + dx
				ys := int(math.Round((sinma*float64(xt)+cosma*float64(yt))+float64(hheight))) + dy

				if xs >= dx && xs < w+dx && ys >= dy && ys < h+dy {
					rgb := img.At(xs, ys)
					//fmt.Println(rgb)
					nrgba.Set(x, y, rgb)
				} else {
					nrgba.Set(x, y, bgSrc)
				}
			}
		}

		sr := nrgba.Bounds()
		r := sr.Sub(sr.Min).Add(image.Point{dx, dy})
		draw.Draw(img, r, nrgba, image.ZP, draw.Src)
	}

	canvasPool.Put(nrgba)
}

// drawLine from min point to max point
func drawLine(m *image.RGBA, x, y, x1, y1 int, c color.Color) {
	if x > x1 || y > y1 {
		x, x1 = x1, x
		y, y1 = y1, y
	}

	rgb := image.NewUniform(c)

	if x == x1 {
		for i := 0; i < y1-y; i++ {
			m.Set(x, y+i, rgb)
		}
	} else if y == y1 {
		for i := 0; i < x1-x; i++ {
			m.Set(x+i, y, rgb)
		}
	} else {
		bresenham(m, x, y, x1, y1, c)
	}
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Bresenham line drawing algorithm
func bresenham(m *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		m.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
