package captchy

import (
	"math"
	"math/rand"
	"time"
)

const letterBytes = "ABDEFHKLMNPRSTUVWXZabdefgikmnopqrstuvwxyz023456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytes(n int) []byte {
	var source = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, source.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// RandomString is a goroutine safe string generator
func RandomString() []byte {
	return randStringBytes(capLen)
}

func devateRandI(devation int) int {
	return rand.Intn(devation<<1) - devation
}

func devateRandF(devation float64) float64 {
	return (rand.Float64() * devation * 2) - devation
}

// randAngle return random radian in -degree ~ degree
func randAngle(degree float64) float64 {
	return (math.Pi / 180) * devateRandF(degree)
}
