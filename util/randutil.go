package util

import (
	"math/rand"
	"time"
)

var (
	randWithSrc *rand.Rand
)

func init() {
	randWithSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt() int {
	return randWithSrc.Int()
}
