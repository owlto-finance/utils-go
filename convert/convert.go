package convert

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func ConvertStringToInt32(value string) int32 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int32(result)
}

func ConvertStringToInt(value string) int {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int(result)
}

func ConvertStringToInt64(value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func ConvertStringToUint64(value string) uint64 {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func ConvertStringToPtrTime(value string) *time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &t
}

func FormatDecimalString(inputStr string, decimalPlaces int) string {
	if inputStr == "" {
		return "0." + strings.Repeat("0", decimalPlaces)
	}

	inputStr = strings.Replace(inputStr, ".", "", 1)

	inputLen := len(inputStr)
	if inputLen <= decimalPlaces {
		return "0." + strings.Repeat("0", decimalPlaces-inputLen) + inputStr
	}

	pointPosition := inputLen - decimalPlaces
	integerPart := inputStr[:pointPosition]
	decimalPart := inputStr[pointPosition:]

	decimalPart = strings.TrimRight(decimalPart, "0")

	if decimalPart == "" {
		return integerPart
	}

	return integerPart + "." + decimalPart
}

// ConvertAndScale 通过缩放10的倍数，scale缩放倍数
func ConvertAndScale(input string, scale int32) string {
	value, _, _ := new(big.Float).SetPrec(236).Parse(input, 10)
	return value.Quo(value, big.NewFloat(math.Pow10(int(scale)))).String()
}
