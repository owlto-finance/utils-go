package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/ethereum/go-ethereum/common"
)

func GetChecksumAddress(address string) (string, error) {
	address = strings.TrimSpace(address)
	prefixLen := 0
	if strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X") {
		prefixLen = 2
	}
	if len(address)-prefixLen == 40 {
		return GetChecksumAddress40(address)
	} else if len(address)-prefixLen == 64 {
		return GetChecksumAddress64(address)
	} else {
		return "", fmt.Errorf("unsupport address len: %s", address)
	}
}

func GetChecksumAddress40(address string) (string, error) {
	address = strings.TrimSpace(address)
	return common.HexToAddress(address).Hex(), nil
}

func GetChecksumAddress64(address string) (string, error) {
	address = strings.TrimSpace(address)
	address = strings.TrimLeft(strings.TrimPrefix(strings.ToLower(address), "0x"), "0")
	if len(address)%2 != 0 {
		address = "0" + address
	}
	if len(address) > 64 {
		return "", errors.New("address too long")
	}
	address64 := strings.Repeat("0", 64-len(address)) + address
	chars := strings.Split(address64, "")
	byteSlice, err := hex.DecodeString(address)
	if err != nil {
		return "", err
	}
	h, err := curve.Curve.StarknetKeccak(byteSlice)
	if err != nil {
		return "", err
	}
	hs := strings.TrimPrefix(h.String(), "0x")
	fmt.Println("1", hs)
	hashed, err := hex.DecodeString(strings.Repeat("0", 64-len(hs)) + hs)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(chars); i += 2 {
		if hashed[i>>1]>>4 >= 8 {
			chars[i] = strings.ToUpper(chars[i])
		}
		if (hashed[i>>1] & 0x0f) >= 8 {
			chars[i+1] = strings.ToUpper(chars[i+1])
		}
	}
	return "0x" + strings.Join(chars, ""), nil
}
