package helpers

import (
	"math/rand"
	"time"
)

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateToken(v int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, v)
	for i := range b {
		b[i] = letterBytes[seededRand.Intn(len(letterBytes))]
	}
	return string(b)
}
