# Package captchy

![build](https://img.shields.io/travis/yujiahaol68/captchy/master.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/yujiahaol68/captchy)](https://goreportcard.com/report/github.com/yujiahaol68/captchy)
[![GoDoc](https://godoc.org/github.com/yujiahaol68/captchy?status.svg)](https://godoc.org/github.com/yujiahaol68/captchy)

```go
import "github.com/yujiahaol68/captchy"
```

captchy implements a abundant CAPTCHAs image generation through flexible config

Split and rotate subImage to against OCRs

Support custom TTF font file or Go builtin fonts

Support export to Base64, PNG, JPEG image

Only std lib so it is clean

**API is still not stable. Not recommend for using in production env !**

## Preview

Generate from very simple to very complex captcha especially for machine, but still friendly to human

![simple-captcha-png](https://github.com/yujiahaol68/captchy/blob/master/example/simple-out.png?raw=true)

![complex-captcha-jpg](https://github.com/yujiahaol68/captchy/blob/master/example/j-out.jpeg?raw=true)

![complex-captcha-png](https://github.com/yujiahaol68/captchy/blob/master/example/p-out.png?raw=true)

## Usage

### Initialize resource

```go
// Load fonts and setup default config
captchy.New(captchy.Default())
```

### Gen solution

```go
rbs := captchy.RandomString()
// Then hash and save rbs into session by yourself
// ..
```

### Gen Image

```go
encoder := captchy.GenerateImg(rbs)
// encoder can encode image into io.Writer, and write it anywhere you prefer.
// It can be local file writer or http response writer
// ..
```

## Reference

[DESIGNING CAPTCHA ALGORITHM: SPLITTING AND ROTATING THE IMAGES AGAINST OCRs](http://cmp.felk.cvut.cz/~cernyad2/TextCaptchaPdf/DESIGNING%20CAPTCHA%20ALGORITHM%20SPLITTING%20AND%20ROTATING.pdf)