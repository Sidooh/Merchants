package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomBool generates a random boolean
func RandomBool() bool {
	return rand.Intn(2) >= 1
}

// RandomIntBetween generates a random integer between min and max
func RandomIntBetween(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomInt generates a random integer between min and max
func RandomInt(length int) int64 {
	max := int64(math.Pow10(length)) - 1
	min := int64(math.Pow10(length - 1))

	return RandomIntBetween(min, max)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@%s.%s", RandomString(6), RandomString(4), RandomString(3))
}

// RandomPhone generates a random phone
func RandomPhone() string {
	prefix := "7"
	if RandomBool() {
		prefix = "1"
	}

	return fmt.Sprintf("%s%s", prefix, strconv.FormatInt(RandomInt(8), 10))
}
