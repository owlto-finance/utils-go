package util

import (
	"math/big"
	"strings"
)

func IsHexStringZero(hexString string) bool {
	// Remove "0x" prefix if it exists
	hexString = strings.TrimSpace(hexString)
	if len(hexString) >= 2 && (hexString[:2] == "0x" || hexString[:2] == "0X") {
		hexString = hexString[2:]
	}

	// Check if all characters in the hex string are '0'
	for _, ch := range hexString {
		if ch != '0' {
			return false
		}
	}
	return true
}

func GetJsonBigInt(itf interface{}) *big.Int {
	switch itf := itf.(type) {
	case float64:
		return big.NewInt(int64(itf))
	case string:
		bi := new(big.Int)
		bi, success := bi.SetString(strings.TrimSpace(itf), 0)
		if success {
			return bi
		} else {
			return new(big.Int)
		}
	default:
		return new(big.Int)
	}
}

func FromUiAmount(amountString string, decimals int) (*big.Int, error) {
	// Convert the amount string to a big.Float
	amountFloat, _, err := new(big.Float).SetPrec(236).Parse(amountString, 10)
	if err != nil {
		return nil, err
	}

	// Scale the amount by 10^(decimals)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountScaled := new(big.Float).SetPrec(236).Mul(amountFloat, new(big.Float).SetPrec(236).SetInt(scale))

	// Convert the scaled amount to a big.Int
	amountBigInt := new(big.Int)
	amountScaled.Int(amountBigInt)

	return amountBigInt, nil
}
