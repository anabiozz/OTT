package utils

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letters = "abcdefghijklmnopqrstuvwxyz"
const (
	letterBits = 5                       // 6 bits to represent a letter index
	letterMask = 1<<letterBits - 1       // All 1-bits, as many as letterIdxBits
	letterMax  = letterMask / letterBits // # of letter indices fitting in 63 bits
)

// Generator ...
func Generator() <-chan string {
	stream := make(chan string, 10)
	go func() {
		for {
			stream <- stringGenerator(5)
			time.Sleep(300 * time.Millisecond)
		}
	}()
	return stream
}

func stringGenerator(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterMax
		}
		if idx := int(cache & letterMax); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterBits
		remain--
	}
	return string(b)
}
