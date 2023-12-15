package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var rd *rand.Rand

func init() {
	rd=rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rd.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rd.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUUID generates a random uuid string
func RandomUUID() uuid.UUID {
	newID,_:= uuid.NewUUID()		
	return newID
}

// RandomQuestion generates a random question
func RandomQuestion() string {
	return RandomString(10)
}

// RandomOptionKey generates a random key
func RandomOptionKey() string {
	return RandomString(6)
}

// RandomOptionValue generates a random value
func RandomOptionValue() string {
	return RandomString(20)
}
