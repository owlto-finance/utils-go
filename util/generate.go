package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Shanghai")

func GenerateLogID() string {
	now := time.Now().In(loc)
	random := rand.New(rand.NewSource(now.UnixNano()))
	timestamp := now.Format("20060102150405")
	randomPart := fmt.Sprintf("%018d", random.Int63())
	if len(randomPart) > 18 {
		randomPart = randomPart[:18]
	}
	return fmt.Sprintf("%s%s", timestamp, randomPart)
}

func GetLogId(ctx context.Context) string {
	logId := ctx.Value("logId")
	if LogId, ok := logId.(string); ok {
		return LogId
	}
	return ""
}

func WithLogIDCtx(ctx context.Context, logId string) context.Context {
	return context.WithValue(ctx, "logId", logId)
}
