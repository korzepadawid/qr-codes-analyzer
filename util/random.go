package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		random := rand.Intn(len(alphabet))
		sb.WriteByte(alphabet[random])
	}
	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomMail() string {
	return fmt.Sprintf("%s@gmail.com", RandomUsername())
}

func RandomInt64(start, end int64) int64 {
	return start + rand.Int63n(end-start+1)
}
