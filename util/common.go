package util

func IsHexStringZero(hexString string) bool {
	// Remove "0x" prefix if it exists
	if len(hexString) > 2 && hexString[:2] == "0x" {
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
