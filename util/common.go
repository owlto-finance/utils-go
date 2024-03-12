package util

import "math/big"

func IsHexStringZero(hexString string) bool {
	// Remove "0x" prefix if it exists
	if len(hexString) > 2 && (hexString[:2] == "0x" || hexString[:2] == "0X") {
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
		bi, success := bi.SetString((itf), 0)
		if success {
			return bi
		} else {
			return new(big.Int)
		}
	default:
		return new(big.Int)
	}
}
