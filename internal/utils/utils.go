package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func HandleError(Err error) bool {
	if Err != nil {
		fmt.Println(Err)
		return true
	}

	return false
}

func RandHexString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
