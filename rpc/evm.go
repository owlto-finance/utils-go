package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/owlto-finance/utils-go/abi/erc20"
	"github.com/owlto-finance/utils-go/loader"
	"github.com/owlto-finance/utils-go/log"
	"github.com/owlto-finance/utils-go/util"
)

type EvmRpc struct {
	chainInfo *loader.ChainInfo
}

func NewEvmRpc(chainInfo *loader.ChainInfo) *EvmRpc {
	return &EvmRpc{
		chainInfo: chainInfo,
	}
}

func (w *EvmRpc) GetClient() *ethclient.Client {
	return w.chainInfo.Client.(*ethclient.Client)
}

func (w *EvmRpc) Client() interface{} {
	return w.chainInfo.Client
}

func (w *EvmRpc) Backend() int32 {
	return 1
}

func (w *EvmRpc) GetAllowance(ctx context.Context, ownerAddr string, tokenAddr string, spenderAddr string) (*big.Int, error) {
	econtract, err := erc20.NewErc20(common.HexToAddress(tokenAddr), w.GetClient())
	if err != nil {
		return nil, err
	}
	allowance, err := econtract.Allowance(nil, common.HexToAddress(ownerAddr), common.HexToAddress(spenderAddr))

	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (w *EvmRpc) GetBalance(ctx context.Context, ownerAddr string, tokenAddr string) (*big.Int, error) {
	ownerAddr = strings.TrimSpace(ownerAddr)
	tokenAddr = strings.TrimSpace(tokenAddr)

	if util.IsHexStringZero(tokenAddr) {
		nativeBalance, err := w.GetClient().BalanceAt(ctx, common.HexToAddress(ownerAddr), nil)
		if err != nil {
			return nil, err
		}
		return nativeBalance, nil
	} else {
		econtract, err := erc20.NewErc20(common.HexToAddress(tokenAddr), w.GetClient())
		if err != nil {
			return nil, err
		}
		balance, err := econtract.BalanceOf(nil, common.HexToAddress(ownerAddr))

		if err != nil {
			return nil, err
		}
		return balance, nil
	}
}

func (w *EvmRpc) IsTxSuccess(ctx context.Context, hash string) (bool, int64, error) {
	receipt, err := w.GetClient().TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return false, 0, err
	}
	if receipt == nil {
		return false, 0, fmt.Errorf("get receipt failed")
	}
	return receipt.Status == ethtypes.ReceiptStatusSuccessful, receipt.BlockNumber.Int64(), nil
}

func (w *EvmRpc) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	blockNumber, err := w.GetClient().BlockNumber(ctx)
	if err != nil {
		log.Errorf("%v get latest block number error %v", w.chainInfo.Name, err)
		return 0, err
	}
	return int64(blockNumber), nil
}
