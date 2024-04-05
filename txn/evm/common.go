package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func EstimateGas(client *ethclient.Client, from string, to string, value *big.Int, data []byte) (uint64, error) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)

	f := common.HexToAddress(from)
	t := common.HexToAddress(to)
	v := value
	if v == nil {
		v = big.NewInt(0)
	}
	if data == nil {
		data = []byte{}
	}
	cm := ethereum.CallMsg{
		From:  f,
		To:    &t,
		Value: v,
		Data:  data,
	}

	gas, err := client.EstimateGas(context.Background(), cm)
	gas = gas * 3 / 2
	return gas, err
}

func TransferBody(client *ethclient.Client, senderAddr string, receiverAddr string, amount *big.Int) ([]byte, error) {
	senderAddr = strings.TrimSpace(senderAddr)
	receiverAddr = strings.TrimSpace(receiverAddr)

	gas, err := EstimateGas(client, senderAddr, receiverAddr, amount, nil)
	if err != nil {
		return nil, err
	}
	return ToBody(receiverAddr, amount, nil, gas)

}

func ToBody(to string, value *big.Int, input []byte, gas uint64) ([]byte, error) {
	to = strings.TrimSpace(to)
	t := common.HexToAddress(to)
	v := value
	if v == nil {
		v = big.NewInt(0)
	}
	m := map[string]interface{}{
		"to":    t.Hex(),
		"gas":   fmt.Sprintf("0x%x", gas),
		"value": fmt.Sprintf("0x%x", v),
	}

	if input != nil {
		m["input"] = fmt.Sprintf("0x%s", hex.EncodeToString(input))
	}

	// Marshal the map to a JSON string
	return json.Marshal(m)
}
