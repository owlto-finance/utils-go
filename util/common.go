package util

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

func FromUiString(amount string, decimals int32) (*big.Int, error) {
	// Convert the amount string to a big.Float
	amountFloat, _, err := new(big.Float).SetPrec(236).Parse(amount, 10)
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

func FromUiFloat(amount float64, decimals int32) *big.Int {
	// Convert the amount string to a big.Float
	amountFloat := new(big.Float).SetPrec(236).SetFloat64(amount)

	// Scale the amount by 10^(decimals)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountScaled := new(big.Float).SetPrec(236).Mul(amountFloat, new(big.Float).SetPrec(236).SetInt(scale))

	// Convert the scaled amount to a big.Int
	amountBigInt := new(big.Int)
	amountScaled.Int(amountBigInt)

	return amountBigInt
}

func StringToUi(amountStr string, decimals int32) (*big.Float, error) {
	// Parse amount string to a big.Int
	amountBigInt, success := new(big.Int).SetString(amountStr, 10)
	if !success {
		return nil, fmt.Errorf("invalid amount string: %s", amountStr)
	}

	// Convert amount to a big.Float
	amountFloat := new(big.Float).SetPrec(236).SetInt(amountBigInt)

	// Scale down the amount by 10^(decimals)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	divider := new(big.Float).SetPrec(236).SetInt(scale)
	amountFloat.Quo(amountFloat, divider)

	return amountFloat, nil
}

func BigIntToUi(amount *big.Int, decimals int32) *big.Float {
	// Convert amount to a big.Float
	amountFloat := new(big.Float).SetPrec(236).SetInt(amount)

	// Scale down the amount by 10^(decimals)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	divider := new(big.Float).SetPrec(236).SetInt(scale)
	amountFloat.Quo(amountFloat, divider)

	return amountFloat
}

func IsEvmAddress(address string, chainID int32) bool {
	return common.IsHexAddress(address) && chainID != 666666666 && chainID != 83797601
}
