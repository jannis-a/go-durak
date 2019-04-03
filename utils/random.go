package utils

import (
	"math/rand"
	"time"
)

const ALPHANUM = "0123456789" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandString(length int) string {
	result := make([]byte, length)

	for i := range result {
		result[i] = ALPHANUM[seed.Intn(len(ALPHANUM))]
	}

	return string(result)
}
