package evm

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/owlto-finance/utils-go/abi/erc20"
)

func Erc20ApproveBody(client *ethclient.Client, senderAddr string, tokenAddr string, spender string, amount *big.Int) ([]byte, error) {
	senderAddr = strings.TrimSpace(senderAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)
	spender = strings.TrimSpace(spender)

	abi, err := erc20.Erc20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("approve", common.HexToAddress(spender), amount)
	if err != nil {
		return nil, err
	}

	gas, err := EstimateGas(client, senderAddr, tokenAddr, nil, data)
	if err != nil {
		return nil, err
	}
	return ToBody(tokenAddr, nil, data, gas)

	// abi, _ := erc20.Erc20MetaData.GetAbi()

	// calldata, err := abi.Pack("approve", common.HexToAddress(spender), amount)
	// if err != nil {
	// 	return nil, err
	// }

	// return calldata, nil
}

func Erc20TransferBody(client *ethclient.Client, senderAddr string, tokenAddr string, receiverAddr string, amount *big.Int) ([]byte, error) {
	senderAddr = strings.TrimSpace(senderAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)
	receiverAddr = strings.TrimSpace(receiverAddr)

	abi, err := erc20.Erc20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("transfer", common.HexToAddress(receiverAddr), amount)
	if err != nil {
		return nil, err
	}

	gas, err := EstimateGas(client, senderAddr, tokenAddr, nil, data)
	if err != nil {
		return nil, err
	}
	return ToBody(tokenAddr, nil, data, gas)
}
