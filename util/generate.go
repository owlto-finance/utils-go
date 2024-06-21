package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateLogID() string {
	now := time.Now()
	random := rand.New(rand.NewSource(now.UnixNano()))
	timestamp := now.Format("20060102150405")
	randomPart := fmt.Sprintf("%018d", random.Int63())
	return fmt.Sprintf("%s%s", timestamp, randomPart)
}
