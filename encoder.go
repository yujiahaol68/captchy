package captchy

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

var _ Encoder = new(CaptchaEncoder)

type Encoder interface {
	ToPNG(io.Writer) error
	ToBase64(io.Writer) error
	ToJPEG(io.Writer, *jpeg.Options) error
}

// CaptchaEncoder encode image.Image into common pic format
type CaptchaEncoder struct {
	img image.Image
}

func (ce *CaptchaEncoder) ToPNG(w io.Writer) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	return png.Encode(buf, ce.img)
}

func (ce *CaptchaEncoder) ToBase64(w io.Writer) error {
	var b bytes.Buffer
	err := png.Encode(&b, ce.img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(base64.StdEncoding.EncodeToString(b.Bytes())))
	return err
}

func (ce *CaptchaEncoder) ToJPEG(w io.Writer, opt *jpeg.Options) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	return jpeg.Encode(buf, ce.img, opt)
}
